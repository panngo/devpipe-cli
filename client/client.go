package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/panngo/devpipe-cli/config"
	"github.com/panngo/devpipe-cli/ws"
)

type IncomingRequest struct {
	ID      string            `json:"id"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type OutgoingResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// Supported HTTP methods that browsers can send
var supportedMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"PATCH":   true,
	"HEAD":    true,
	"OPTIONS": true,
	"TRACE":   true,
	"CONNECT": true,
}

// Methods that typically have a body
var methodsWithBody = map[string]bool{
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": false, // DELETE can have body but often doesn't
}

func ParseFlags() string {
	port := flag.String("port", "3000", "Local port to forward to")
	clearConfig := flag.Bool("clear-config", false, "Clear saved tunnel configuration")
	flag.Parse()
	
	// Handle clear config flag
	if *clearConfig {
		configManager := config.NewConfigManager()
		if err := configManager.ClearTunnelConfig(); err != nil {
			log.Printf("❌ Failed to clear tunnel configuration: %v", err)
		} else {
			log.Println("🗑️  Tunnel configuration cleared successfully")
		}
	}
	
	return *port
}

func ListenAndServe(conn *ws.SafeConn, port string) {
	serverUrl := "wss://devpipe.cloud/ws"
	
	// Store the initial tunnel ID and UUID
	tunnelID := conn.GetTunnelID()
	uuid := conn.GetUUID()
	
	log.Printf("🔗 Connected with tunnel: %s", tunnelID)
	if uuid != "" {
		log.Printf("🔑 UUID: %s", uuid)
	}
	
	// Configure heartbeat
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()
	
	// Channel for communication between goroutines
	errorChan := make(chan error, 1)
	
	// Goroutine for heartbeat
	go func() {
		for range heartbeatTicker.C {
			if err := conn.WriteJSON(map[string]string{"action": "ping"}); err != nil {
				errorChan <- err
				return
			}
		}
	}()
	
	for {
		select {
		case err := <-errorChan:
			log.Println("❌ Heartbeat error:", err)
			goto reconnect
		default:
			// Configure timeout for reading
			conn.SetReadDeadline(time.Now().Add(35 * time.Second))
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("❌ WebSocket read error:", err)
				goto reconnect
			}
			
			// Reset deadline after successful read
			conn.SetReadDeadline(time.Time{})
			
			var req IncomingRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				log.Println("❌ JSON unmarshal error:", err)
				continue
			}

			go handleRequest(conn, req, port)
		}
		continue
	
	reconnect:
		log.Println("🔄 Attempting to reconnect...")
		
		// Stop heartbeat temporarily
		heartbeatTicker.Stop()
		
		// Try to reconnect
		newConn, newTunnelID := reconnect(serverUrl, port, tunnelID, uuid)
		if newConn == nil {
			log.Println("❌ Failed to reconnect, exiting...")
			return
		}
		
		conn = newConn
		tunnelID = newTunnelID
		uuid = conn.GetUUID()
		
		if tunnelID == conn.GetTunnelID() {
			log.Printf("✅ Successfully reconnected with same tunnel: %s", tunnelID)
		} else {
			log.Printf("⚠️  Reconnected with new tunnel: %s (was: %s)", tunnelID, conn.GetTunnelID())
		}
		
		// Restart heartbeat
		heartbeatTicker = time.NewTicker(30 * time.Second)
		go func() {
			for range heartbeatTicker.C {
				if err := conn.WriteJSON(map[string]string{"action": "ping"}); err != nil {
					errorChan <- err
					return
				}
			}
		}()
	}
}

func reconnect(serverUrl, port, previousTunnelID, previousUUID string) (*ws.SafeConn, string) {
	maxRetries := 5
	retryDelay := time.Second * 2
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 Reconnection attempt %d/%d...", attempt, maxRetries)
		
		var conn *ws.SafeConn
		var tunnelID string
		var err error
		
		// If we have a previous UUID, try secure reconnection first
		if previousUUID != "" && attempt == 1 {
			log.Printf("🔐 Attempting secure reconnection with UUID: %s", previousUUID)
			conn, tunnelID, err = ws.ConnectAndReconnect(serverUrl, port, previousTunnelID)
			if err == nil {
				log.Printf("✅ Secure reconnection successful")
				return conn, tunnelID
			} else {
				log.Printf("❌ Secure reconnection failed: %v", err)
				// Clear invalid configuration
				configManager := config.NewConfigManager()
				if clearErr := configManager.ClearTunnelConfig(); clearErr != nil {
					log.Printf("⚠️  Warning: Could not clear invalid config: %v", clearErr)
				} else {
					log.Printf("🗑️  Cleared invalid tunnel configuration")
				}
			}
		}
		
		// Fallback to new registration
		log.Println("🆕 Attempting new registration...")
		conn, tunnelID, err = ws.ConnectAndRegisterWithRetry(serverUrl, port)
		
		if err == nil {
			// If we had a previous tunnel ID and the new one is different, log this
			if previousTunnelID != "" && previousTunnelID != tunnelID {
				log.Printf("⚠️  New tunnel ID assigned: %s (was: %s)", tunnelID, previousTunnelID)
			} else if previousTunnelID == tunnelID {
				log.Printf("✅ Successfully reconnected with same tunnel: %s", tunnelID)
			}
			return conn, tunnelID
		}
		
		log.Printf("❌ Reconnection attempt %d failed: %v", attempt, err)
		
		if attempt < maxRetries {
			log.Printf("⏳ Waiting %v before next attempt...", retryDelay)
			time.Sleep(retryDelay)
			// Increase delay exponentially
			retryDelay = time.Duration(float64(retryDelay) * 1.5)
		}
	}
	
	return nil, ""
}

func handleRequest(conn *ws.SafeConn, req IncomingRequest, port string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Panic in request handler: %v", r)
		}
	}()
	
	// Validate HTTP method
	if !isValidHTTPMethod(req.Method) {
		log.Printf("❌ Unsupported HTTP method: %s", req.Method)
		sendErrorResponse(conn, req.ID, fmt.Sprintf("Unsupported HTTP method: %s", req.Method), 405)
		return
	}
	
	// Validate request path
	if req.Path == "" {
		log.Printf("❌ Empty request path")
		sendErrorResponse(conn, req.ID, "Empty request path", 400)
		return
	}
	
	// Handle special methods
	if req.Method == "OPTIONS" {
		handleOptionsRequest(conn, req, port)
		return
	}
	
	url := "http://localhost:" + port + req.Path
	
	// Create request with appropriate body handling
	var httpReq *http.Request
	var err error
	
	if shouldHaveBody(req.Method) && req.Body != "" {
		httpReq, err = http.NewRequest(req.Method, url, strings.NewReader(req.Body))
	} else {
		httpReq, err = http.NewRequest(req.Method, url, nil)
	}
	
	if err != nil {
		log.Printf("❌ Error creating local request for %s: %v", url, err)
		sendErrorResponse(conn, req.ID, "Failed to create request", 500)
		return
	}
	
	// PROXY MODE: Copy ALL headers exactly as received (transparent proxy)
	for k, v := range req.Headers {
		// Only skip headers that Go's HTTP client manages automatically
		if strings.ToLower(k) == "host" {
			// Skip Host header, it will be defined automatically by Go
			continue
		}
		if strings.ToLower(k) == "connection" {
			// Skip Connection header, it will be managed by Go
			continue
		}
		if strings.ToLower(k) == "transfer-encoding" {
			// Skip Transfer-Encoding, let Go handle it
			continue
		}
		// Copy all other headers exactly as received
		httpReq.Header.Set(k, v)
	}
	
	// Log the request for debugging
	log.Printf("🌐 HTTP %s %s", req.Method, req.Path)
	if req.Body != "" {
		log.Printf("📦 Request body length: %d bytes", len(req.Body))
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("❌ Request failed: %v", err)
		sendErrorResponse(conn, req.ID, "Request failed", 502)
		return
	}
	defer resp.Body.Close()

	// Handle HEAD requests specially (no body)
	if req.Method == "HEAD" {
		handleHeadResponse(conn, req, resp)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Error reading response body: %v", err)
		sendErrorResponse(conn, req.ID, "Failed to read response", 500)
		return
	}

	response := OutgoingResponse{
		ID:      req.ID,
		Status:  resp.StatusCode,
		Headers: map[string]string{},
		Body:    string(body),
	}
	
	// PROXY MODE: Copy ALL response headers exactly as received (transparent proxy)
	for k, v := range resp.Header {
		if len(v) > 0 {
			// For headers that can have multiple values, join with comma
			if len(v) > 1 {
				response.Headers[k] = strings.Join(v, ", ")
			} else {
				response.Headers[k] = v[0]
			}
		}
	}
	
	// Ensure Content-Length is calculated correctly
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	
	fmt.Printf("%-6s %-20s %d OK\n", req.Method, req.Path, resp.StatusCode)
	
	// Use a mutex or channel to ensure thread-safety in writing
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("❌ Error sending response: %v", err)
	}
}

// isValidHTTPMethod checks if the method is supported
func isValidHTTPMethod(method string) bool {
	return supportedMethods[strings.ToUpper(method)]
}

// shouldHaveBody determines if a method typically has a body
func shouldHaveBody(method string) bool {
	return methodsWithBody[strings.ToUpper(method)]
}

// handleOptionsRequest handles OPTIONS requests (CORS preflight)
func handleOptionsRequest(conn *ws.SafeConn, req IncomingRequest, port string) {
	response := OutgoingResponse{
		ID:     req.ID,
		Status: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With, Accept, Origin, User-Agent, Referer",
			"Access-Control-Max-Age":           "86400",
			"Access-Control-Allow-Credentials": "true",
			"Content-Length":                   "0",
		},
		Body: "",
	}
	
	log.Printf("🌐 HTTP OPTIONS %s (CORS preflight)", req.Path)
	fmt.Printf("%-6s %-20s %d OK\n", "OPTIONS", req.Path, 200)
	
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("❌ Error sending OPTIONS response: %v", err)
	}
}

// handleHeadResponse handles HEAD requests (no body)
func handleHeadResponse(conn *ws.SafeConn, req IncomingRequest, resp *http.Response) {
	response := OutgoingResponse{
		ID:     req.ID,
		Status: resp.StatusCode,
		Headers: map[string]string{},
		Body:    "", // HEAD requests should not have a body
	}
	
	// Copy all headers from the response
	for k, v := range resp.Header {
		if len(v) > 0 {
			if len(v) > 1 {
				response.Headers[k] = strings.Join(v, ", ")
			} else {
				response.Headers[k] = v[0]
			}
		}
	}
	
	// Ensure Content-Length is set to 0 for HEAD requests
	response.Headers["Content-Length"] = "0"
	
	log.Printf("🌐 HTTP HEAD %s (no body)", req.Path)
	fmt.Printf("%-6s %-20s %d OK\n", "HEAD", req.Path, resp.StatusCode)
	
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("❌ Error sending HEAD response: %v", err)
	}
}

func sendErrorResponse(conn *ws.SafeConn, reqID, message string, status int) {
	response := OutgoingResponse{
		ID:     reqID,
		Status: status,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: message,
	}
	
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("❌ Error sending error response: %v", err)
	}
}