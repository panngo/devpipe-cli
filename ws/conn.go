package ws

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/panngo/devpipe-cli/config"
)

type SafeConn struct {
	*websocket.Conn
	TunnelID    string
	UUID        string
	SecurityKey string
	URL         string
	writeMutex  sync.Mutex
}

// WriteJSON thread-safe wrapper
func (s *SafeConn) WriteJSON(v interface{}) error {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()
	return s.Conn.WriteJSON(v)
}

// RegistrationResponse represents the server response for registration
type RegistrationResponse struct {
	Tunnel      string `json:"tunnel"`
	UUID        string `json:"uuid"`
	SecurityKey string `json:"key"`
	URL         string `json:"url"`
	Error       string `json:"error,omitempty"`
}

func ConnectAndRegister(serverUrl, port string) (*SafeConn, string) {
	configManager := config.NewConfigManager()

	// Try to load existing tunnel configuration
	existingConfig, err := configManager.LoadTunnelConfig()
	if err != nil {
		log.Printf("⚠️  Warning: Could not load existing config: %v", err)
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(serverUrl, nil)
	if err != nil {
		log.Fatalf("❌ connection error: %v", err)
	}

	safeConn := &SafeConn{Conn: conn}

	// Prepare registration message
	registration := map[string]string{
		"action": "register",
		"port":   port,
	}

	// If we have existing config with UUID and security key, try secure reconnection
	if existingConfig != nil && existingConfig.UUID != "" && existingConfig.SecurityKey != "" {
		log.Printf("🔐 Attempting secure reconnection with UUID: %s", existingConfig.UUID)
		registration["uuid"] = existingConfig.UUID
		registration["key"] = existingConfig.SecurityKey
		safeConn.UUID = existingConfig.UUID
		safeConn.SecurityKey = existingConfig.SecurityKey
	} else {
		log.Println("🆕 Creating new secure connection")
	}

	if err := safeConn.WriteJSON(registration); err != nil {
		log.Fatalf("❌ Failed to send registration: %v", err)
	}

	var response RegistrationResponse
	if err := conn.ReadJSON(&response); err != nil {
		log.Fatalf("❌ Failed to read registration response: %v", err)
	}

	// Check for server error
	if response.Error != "" {
		log.Fatalf("❌ Server error: %s", response.Error)
	}

	// Update connection with new tunnel info
	safeConn.TunnelID = response.Tunnel
	safeConn.UUID = response.UUID
	safeConn.SecurityKey = response.SecurityKey
	safeConn.URL = response.URL

	// Save the new configuration
	newConfig := config.TunnelConfig{
		UUID:        response.UUID,
		SecurityKey: response.SecurityKey,
		TunnelID:    response.Tunnel,
		Port:        port,
	}

	if err := configManager.SaveTunnelConfig(newConfig); err != nil {
		log.Printf("⚠️  Warning: Could not save tunnel config: %v", err)
	} else {
		log.Printf("💾 Tunnel configuration saved for secure reconnection")
	}

	return safeConn, response.Tunnel
}

func ConnectAndRegisterWithRetry(serverUrl, port string) (*SafeConn, string, error) {
	configManager := config.NewConfigManager()

	// Try to load existing tunnel configuration
	existingConfig, err := configManager.LoadTunnelConfig()
	if err != nil {
		log.Printf("⚠️  Warning: Could not load existing config: %v", err)
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(serverUrl, nil)
	if err != nil {
		return nil, "", err
	}

	safeConn := &SafeConn{Conn: conn}

	// Prepare registration message
	registration := map[string]string{
		"action": "register",
		"port":   port,
	}

	// If we have existing config with UUID and security key, try secure reconnection
	if existingConfig != nil && existingConfig.UUID != "" && existingConfig.SecurityKey != "" {
		log.Printf("🔐 Attempting secure reconnection with UUID: %s", existingConfig.UUID)
		registration["uuid"] = existingConfig.UUID
		registration["key"] = existingConfig.SecurityKey
		safeConn.UUID = existingConfig.UUID
		safeConn.SecurityKey = existingConfig.SecurityKey
	} else {
		log.Println("🆕 Creating new secure connection")
	}

	if err := safeConn.WriteJSON(registration); err != nil {
		conn.Close()
		return nil, "", err
	}

	var response RegistrationResponse
	if err := conn.ReadJSON(&response); err != nil {
		conn.Close()
		return nil, "", err
	}

	// Check for server error
	if response.Error != "" {
		conn.Close()
		return nil, "", fmt.Errorf("server error: %s", response.Error)
	}

	// Update connection with new tunnel info
	safeConn.TunnelID = response.Tunnel
	safeConn.UUID = response.UUID
	safeConn.SecurityKey = response.SecurityKey
	safeConn.URL = response.URL

	// Save the new configuration
	newConfig := config.TunnelConfig{
		UUID:        response.UUID,
		SecurityKey: response.SecurityKey,
		TunnelID:    response.Tunnel,
		Port:        port,
	}

	if err := configManager.SaveTunnelConfig(newConfig); err != nil {
		log.Printf("⚠️  Warning: Could not save tunnel config: %v", err)
	} else {
		log.Printf("💾 Tunnel configuration saved for secure reconnection")
	}

	return safeConn, response.Tunnel, nil
}

// ConnectAndReconnect attempts to reconnect with a specific tunnel ID using secure reconnection
func ConnectAndReconnect(serverUrl, port, tunnelID string) (*SafeConn, string, error) {
	configManager := config.NewConfigManager()

	// Load existing tunnel configuration
	existingConfig, err := configManager.LoadTunnelConfig()
	if err != nil {
		return nil, "", fmt.Errorf("failed to load tunnel config: %w", err)
	}

	if existingConfig == nil || existingConfig.UUID == "" || existingConfig.SecurityKey == "" {
		return nil, "", fmt.Errorf("no valid tunnel configuration found for reconnection")
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(serverUrl, nil)
	if err != nil {
		return nil, "", err
	}

	safeConn := &SafeConn{Conn: conn}

	// Attempt secure reconnection with UUID and security key
	log.Printf("🔐 Attempting secure reconnection with UUID: %s", existingConfig.UUID)
	registration := map[string]string{
		"action": "register",
		"port":   port,
		"uuid":   existingConfig.UUID,
		"key":    existingConfig.SecurityKey,
	}

	if err := safeConn.WriteJSON(registration); err != nil {
		conn.Close()
		return nil, "", err
	}

	var response RegistrationResponse
	if err := conn.ReadJSON(&response); err != nil {
		conn.Close()
		return nil, "", err
	}

	// Check for server error
	if response.Error != "" {
		conn.Close()
		return nil, "", fmt.Errorf("server error: %s", response.Error)
	}

	// Verify we got the same tunnel ID back
	if response.Tunnel != tunnelID {
		log.Printf("⚠️  Server returned different tunnel ID: %s (expected: %s)", response.Tunnel, tunnelID)
	}

	// Update connection with tunnel info
	safeConn.TunnelID = response.Tunnel
	safeConn.UUID = response.UUID
	safeConn.SecurityKey = response.SecurityKey
	safeConn.URL = response.URL

	return safeConn, response.Tunnel, nil
}

// GetTunnelID returns the tunnel ID of the connection
func (s *SafeConn) GetTunnelID() string {
	return s.TunnelID
}

// GetUUID returns the UUID of the connection
func (s *SafeConn) GetUUID() string {
	return s.UUID
}

// GetSecurityKey returns the security key of the connection
func (s *SafeConn) GetSecurityKey() string {
	return s.SecurityKey
}

// GetURL returns the tunnel URL of the connection
func (s *SafeConn) GetURL() string {
	return s.URL
}
