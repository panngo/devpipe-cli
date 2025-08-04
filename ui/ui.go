package ui

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintBanner(port, tunnelID, url string) {
	clearConsole()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Println(cyan("@devpipe"))
	fmt.Println()
	fmt.Printf("%-15s %s\n", "Tunnel Status", green("online"))
	fmt.Printf("%-15s %s\n", "Version", "custom-devpipe")
	fmt.Printf("%-15s %s\n", "Forwarding", fmt.Sprintf("%s -> localhost:%s", yellow(url), port))
	fmt.Printf("%-15s %s\n", "Security", blue("🔐 Secure Reconnection Enabled"))
	fmt.Println()
	fmt.Println("HTTP Requests")
	fmt.Printf("%-6s %-20s %-6s\n", "METHOD", "PATH", "STATUS")
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}
