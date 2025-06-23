package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

func ParseFlags() string {
	port := flag.String("port", "3000", "Local port to forward to")
	flag.Parse()
	return *port
}

func ListenAndServe(conn *ws.SafeConn, port string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WebSocket read error:", err)
			return
		}

		var req IncomingRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("❌ JSON unmarshal error:", err)
			continue
		}

		go handleRequest(conn, req, port)
	}
}

func handleRequest(conn *ws.SafeConn, req IncomingRequest, port string) {
	url := "http://localhost:" + port + req.Path
	httpReq, err := http.NewRequest(req.Method, url, strings.NewReader(req.Body))
	if err != nil {
		log.Printf("❌ Error creating local request for %s: %v", url, err)
		return
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("❌ Request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	response := OutgoingResponse{
		ID:      req.ID,
		Status:  resp.StatusCode,
		Headers: map[string]string{},
		Body:    string(body),
	}
	for k, v := range resp.Header {
		if len(v) > 0 {
			response.Headers[k] = v[0]
		}
	}

	fmt.Printf("%-6s %-20s %d OK\n", req.Method, req.Path, resp.StatusCode)
	conn.WriteJSON(response)
}