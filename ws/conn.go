package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type SafeConn struct {
	*websocket.Conn
}

func ConnectAndRegister(serverUrl, port string) (*SafeConn, string) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(serverUrl, nil)
	if err != nil {
		log.Fatalf("❌ connection error: %v", err)
	}

	conn.WriteJSON(map[string]string{
		"action": "register",
		"port":   port,
	})

	var ack struct {
		Tunnel string `json:"tunnel"`
	}
	if err := conn.ReadJSON(&ack); err != nil {
		log.Fatalf("❌ Failed to register: %v", err)
	}

	return &SafeConn{conn}, ack.Tunnel
}