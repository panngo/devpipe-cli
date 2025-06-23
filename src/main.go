package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
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

func main() {
	serverUrl := "wss://devpipe.cloud/ws"

	port := flag.String("port", "3000", "Local port to forward to")
	flag.Parse()

	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(serverUrl, nil)
	if err != nil {
		log.Fatal("❌ connection error: ", err)
	}
	defer conn.Close()

	// Register tunnel
	conn.WriteJSON(map[string]string{
		"action": "register",
		"port":   *port,
	})

	// Read the tunnel ID from the server
	var registrationAck struct {
		Tunnel string `json:"tunnel"`
	}
	if err := conn.ReadJSON(&registrationAck); err != nil {
		log.Fatal("❌ Failed to read tunnel ID from server: ", err)
	}
	tunnelId := registrationAck.Tunnel

	// Print ngrok-style banner
	printNgrokStyleBanner(*port, tunnelId)

	// Log HTTP requests
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WebSocket read error: ", err)
			return
		}

		var req IncomingRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("❌ JSON unmarshal error: ", err)
			continue
		}

		go handleRequest(conn, req, *port)
	}
}

func handleRequest(conn *websocket.Conn, req IncomingRequest, port string) {
	url := "http://localhost:" + port + req.Path
	httpReq, err := http.NewRequest(req.Method, url, strings.NewReader(req.Body))
	if err != nil {
		log.Printf("❌ Error creating local request for %s: %v", url, err)
		return
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("❌ Error making local request to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	respMsg := OutgoingResponse{
		ID:      req.ID,
		Status:  resp.StatusCode,
		Headers: map[string]string{},
		Body:    string(bodyBytes),
	}

	for k, v := range resp.Header {
		if len(v) > 0 {
			respMsg.Headers[k] = v[0]
		}
	}

	// Print request log
	fmt.Printf("%-6s %-20s %d OK\n", req.Method, req.Path, resp.StatusCode)

	conn.WriteJSON(respMsg)
}

func printNgrokStyleBanner(port, tunnelId string) {
	clearConsole()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(cyan("@devpipe"))
	fmt.Println()
	fmt.Printf("%-15s %s\n", "Tunnel Status", green("online"))
	fmt.Printf("%-15s %s\n", "Version", "custom-devpipe")
	fmt.Printf("%-15s %s\n", "Forwarding", fmt.Sprintf("%s -> localhost:%s", yellow("https://"+tunnelId+".devpipe.cloud"), port))
	fmt.Println()
	fmt.Println("HTTP Requests")
	fmt.Printf("%-6s %-20s %-6s\n", "METHOD", "PATH", "STATUS")
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}
