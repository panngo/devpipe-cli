// main.go
package main

import (
	"github.com/panngo/devpipe-cli/client"
	"github.com/panngo/devpipe-cli/ui"
	"github.com/panngo/devpipe-cli/ws"
)

func main() {
	port := client.ParseFlags()

	conn, tunnelID := ws.ConnectAndRegister("wss://devpipe.cloud/ws", port)
	defer conn.Close()

	ui.PrintBanner(port, tunnelID)
	client.ListenAndServe(conn, port)
}