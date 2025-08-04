package interfaces

import (
	"time"

	"github.com/panngo/devpipe-cli/logger"
)

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection interface {
	WriteJSON(v interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	Close() error
	GetTunnelID() string
	GetUUID() string
	GetSecurityKey() string
}

// TunnelManager manages tunnel connections
type TunnelManager interface {
	Register(tunnelID string, conn WebSocketConnection, port string, key string) (string, bool)
	ValidateKey(uuid string, key string) bool
	Get(tunnelID string) Tunnel
	GetActiveTunnelCount() int
}

// Tunnel represents a tunnel connection
type Tunnel interface {
	GetID() string
	GetPort() string
	GetLastSeen() time.Time
	GetConnection() WebSocketConnection
	IsActive() bool
}

// ConfigManager manages configuration
type ConfigManager interface {
	LoadTunnelConfig() (*TunnelConfig, error)
	SaveTunnelConfig(config TunnelConfig) error
	ClearTunnelConfig() error
}

// TunnelConfig represents tunnel configuration
type TunnelConfig struct {
	UUID        string `json:"uuid"`
	SecurityKey string `json:"security_key"`
	TunnelID    string `json:"tunnel_id"`
	Port        string `json:"port"`
}

// HTTPClient represents an HTTP client
type HTTPClient interface {
	Do(req *HTTPRequest) (*HTTPResponse, error)
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	Status  int
	Headers map[string]string
	Body    string
}

// ReconnectionStrategy defines reconnection behavior
type ReconnectionStrategy interface {
	Reconnect(urls []string, port string, tunnelID string, uuid string) (WebSocketConnection, string, error)
}

// HeartbeatService manages heartbeat functionality
type HeartbeatService interface {
	Start(conn WebSocketConnection, interval time.Duration) error
	Stop() error
	IsActive() bool
}

// RequestHandler handles incoming requests
type RequestHandler interface {
	HandleRequest(req *IncomingRequest, port string) (*OutgoingResponse, error)
}

// IncomingRequest represents an incoming request from the server
type IncomingRequest struct {
	ID      string            `json:"id"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// OutgoingResponse represents an outgoing response to the server
type OutgoingResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// Application represents the main application interface
type Application interface {
	Start() error
	Stop() error
	IsRunning() bool
}

// ServiceLocator provides dependency injection
type ServiceLocator interface {
	GetLogger() logger.Logger
	GetTunnelManager() TunnelManager
	GetConfigManager() ConfigManager
	GetHTTPClient() HTTPClient
	GetReconnectionStrategy() ReconnectionStrategy
	GetHeartbeatService() HeartbeatService
	GetRequestHandler() RequestHandler
}
