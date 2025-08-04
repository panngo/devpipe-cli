package services

import (
	"sync"
	"time"

	"github.com/panngo/devpipe-cli/interfaces"
	"github.com/panngo/devpipe-cli/logger"
)

// HeartbeatService manages heartbeat functionality
type HeartbeatService struct {
	logger   logger.Logger
	ticker   *time.Ticker
	stopChan chan struct{}
	isActive bool
	mu       sync.RWMutex
}

// NewHeartbeatService creates a new heartbeat service
func NewHeartbeatService(logger logger.Logger) *HeartbeatService {
	return &HeartbeatService{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

// Start starts the heartbeat service
func (h *HeartbeatService) Start(conn interfaces.WebSocketConnection, interval time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.isActive {
		return nil
	}

	h.ticker = time.NewTicker(interval)
	h.isActive = true

	go h.run(conn)

	h.logger.Debug("Heartbeat service started with interval: %v", interval)
	return nil
}

// Stop stops the heartbeat service
func (h *HeartbeatService) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isActive {
		return nil
	}

	if h.ticker != nil {
		h.ticker.Stop()
	}

	close(h.stopChan)
	h.isActive = false

	h.logger.Debug("Heartbeat service stopped")
	return nil
}

// IsActive returns whether the heartbeat service is active
func (h *HeartbeatService) IsActive() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.isActive
}

// run runs the heartbeat loop
func (h *HeartbeatService) run(conn interfaces.WebSocketConnection) {
	for {
		select {
		case <-h.ticker.C:
			if err := conn.WriteJSON(map[string]string{"action": "ping"}); err != nil {
				h.logger.Error("Failed to send heartbeat: %v", err)
				return
			}
			h.logger.Debug("Heartbeat sent")
		case <-h.stopChan:
			return
		}
	}
}
