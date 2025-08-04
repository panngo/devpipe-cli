package services

import (
	"fmt"
	"time"

	"github.com/panngo/devpipe-cli/interfaces"
	"github.com/panngo/devpipe-cli/logger"
	"github.com/panngo/devpipe-cli/ws"
)

// ReconnectionStrategy implements reconnection logic
type ReconnectionStrategy struct {
	logger        logger.Logger
	configManager interfaces.ConfigManager
	maxRetries    int
	baseDelay     time.Duration
}

// NewReconnectionStrategy creates a new reconnection strategy
func NewReconnectionStrategy(logger logger.Logger, configManager interfaces.ConfigManager) *ReconnectionStrategy {
	return &ReconnectionStrategy{
		logger:        logger,
		configManager: configManager,
		maxRetries:    5,
		baseDelay:     2 * time.Second,
	}
}

// Reconnect attempts to reconnect using multiple strategies
func (r *ReconnectionStrategy) Reconnect(urls []string, port string, tunnelID string, uuid string) (interfaces.WebSocketConnection, string, error) {
	retryDelay := r.baseDelay

	for attempt := 1; attempt <= r.maxRetries; attempt++ {
		r.logger.Debug("Reconnection attempt %d/%d...", attempt, r.maxRetries)

		// Try each URL
		for _, url := range urls {
			r.logger.Debug("Trying to reconnect to: %s", url)

			conn, newTunnelID, err := r.tryReconnect(url, port, tunnelID, uuid)
			if err == nil {
				r.logger.Info("Successfully reconnected to: %s", url)
				return conn, newTunnelID, nil
			}

			r.logger.Debug("Failed to reconnect to %s: %v", url, err)
		}

		if attempt < r.maxRetries {
			r.logger.Debug("Waiting %v before next attempt...", retryDelay)
			time.Sleep(retryDelay)
			retryDelay = time.Duration(float64(retryDelay) * 1.5) // Exponential backoff
		}
	}

	return nil, "", fmt.Errorf("failed to reconnect after %d attempts", r.maxRetries)
}

// tryReconnect attempts to reconnect to a specific URL
func (r *ReconnectionStrategy) tryReconnect(url, port, tunnelID, uuid string) (interfaces.WebSocketConnection, string, error) {
	// Load existing configuration
	existingConfig, err := r.configManager.LoadTunnelConfig()
	if err != nil {
		r.logger.Debug("Could not load existing config: %v", err)
	}

	// Try secure reconnection first if we have valid credentials
	if existingConfig != nil && existingConfig.UUID != "" && existingConfig.SecurityKey != "" {
		r.logger.Debug("Attempting secure reconnection with UUID: %s", existingConfig.UUID)

		conn, newTunnelID, err := ws.ConnectAndReconnect(url, port, tunnelID)
		if err == nil {
			return conn, newTunnelID, nil
		}

		r.logger.Debug("Secure reconnection failed: %v", err)

		// Clear invalid configuration
		if clearErr := r.configManager.ClearTunnelConfig(); clearErr != nil {
			r.logger.Debug("Could not clear invalid config: %v", clearErr)
		}
	}

	// Fallback to new registration
	r.logger.Debug("Attempting new registration...")
	conn, newTunnelID, err := ws.ConnectAndRegisterWithRetry(url, port)
	if err != nil {
		return nil, "", err
	}

	return conn, newTunnelID, nil
}

// SetMaxRetries sets the maximum number of retry attempts
func (r *ReconnectionStrategy) SetMaxRetries(maxRetries int) {
	r.maxRetries = maxRetries
}

// SetBaseDelay sets the base delay between retries
func (r *ReconnectionStrategy) SetBaseDelay(baseDelay time.Duration) {
	r.baseDelay = baseDelay
}
