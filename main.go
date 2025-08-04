// main.go
package main

import (
	"log"

	"github.com/panngo/devpipe-cli/client"
	"github.com/panngo/devpipe-cli/ui"
	"github.com/panngo/devpipe-cli/ws"
)

func main() {
	port := client.ParseFlags()

	// Try different server URLs
	serverURLs := []string{
		"ws://localhost:3000/ws",  // Local development
		// "wss://api.devpipe.cloud/ws",
		// "wss://ws-tunnel.devpipe.cloud/ws",  // Fallback
		// "wss://tunnel.devpipe.cloud/ws",     // Fallback
		// "ws://tunnel.devpipe.cloud/ws",
		// "wss://devpipe.cloud/ws",
		// "ws://devpipe.cloud/ws",
	}
	
	var conn *ws.SafeConn
	var tunnelID string
	var lastErr error
	
	for _, url := range serverURLs {
		log.Printf("🔌 Tentando conectar em: %s", url)
		conn, tunnelID, lastErr = ws.ConnectAndRegisterWithRetry(url, port)
		if lastErr == nil {
			log.Printf("✅ Conectado com sucesso em: %s", url)
			break
		}
		log.Printf("❌ Falha ao conectar em %s: %v", url, lastErr)
	}
	
	if lastErr != nil {
		log.Fatalf("❌ Falha ao conectar em todos os servidores: %v", lastErr)
	}
	defer conn.Close()

	ui.PrintBanner(port, tunnelID)
	ui.PrintSecureReconnectionInfo(conn.GetUUID())
	client.ListenAndServe(conn, port)
}