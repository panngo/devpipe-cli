package services

import (
	"fmt"
	"time"

	"github.com/panngo/devpipe-cli/interfaces"
	"github.com/panngo/devpipe-cli/logger"
)

// RequestHandler handles incoming HTTP requests
type RequestHandler struct {
	logger     logger.Logger
	httpClient interfaces.HTTPClient
	httpLogger *logger.HTTPRequestLogger
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(log logger.Logger, httpClient interfaces.HTTPClient) *RequestHandler {
	return &RequestHandler{
		logger:     log,
		httpClient: httpClient,
		httpLogger: logger.NewHTTPRequestLogger(log),
	}
}

// HandleRequest handles an incoming request
func (r *RequestHandler) HandleRequest(req *interfaces.IncomingRequest, port string) (*interfaces.OutgoingResponse, error) {
	start := time.Now()

	// Validate request
	if err := r.validateRequest(req); err != nil {
		return r.createErrorResponse(req.ID, err.Error(), 400), nil
	}

	// Handle special methods
	if req.Method == "OPTIONS" {
		return r.handleOptionsRequest(req), nil
	}

	// Create HTTP request
	httpReq, err := r.createHTTPRequest(req, port)
	if err != nil {
		return r.createErrorResponse(req.ID, "Failed to create request", 500), nil
	}

	// Execute request
	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return r.createErrorResponse(req.ID, "Failed to execute request", 500), nil
	}

	// Log request
	duration := time.Since(start)
	r.httpLogger.LogRequest(req.Method, req.Path, fmt.Sprintf("%d", resp.Status), duration)

	// Create response
	response := &interfaces.OutgoingResponse{
		ID:      req.ID,
		Status:  resp.Status,
		Headers: resp.Headers,
		Body:    resp.Body,
	}

	return response, nil
}

// validateRequest validates the incoming request
func (r *RequestHandler) validateRequest(req *interfaces.IncomingRequest) error {
	if !r.isValidHTTPMethod(req.Method) {
		return fmt.Errorf("unsupported HTTP method: %s", req.Method)
	}

	if req.Path == "" {
		return fmt.Errorf("empty request path")
	}

	return nil
}

// createHTTPRequest creates an HTTP request
func (r *RequestHandler) createHTTPRequest(req *interfaces.IncomingRequest, port string) (*interfaces.HTTPRequest, error) {
	url := fmt.Sprintf("http://localhost:%s%s", port, req.Path)

	httpReq := &interfaces.HTTPRequest{
		Method:  req.Method,
		URL:     url,
		Headers: req.Headers,
		Body:    req.Body,
	}

	return httpReq, nil
}

// handleOptionsRequest handles OPTIONS requests
func (r *RequestHandler) handleOptionsRequest(req *interfaces.IncomingRequest) *interfaces.OutgoingResponse {
	headers := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Max-Age":           "86400",
	}

	return &interfaces.OutgoingResponse{
		ID:      req.ID,
		Status:  200,
		Headers: headers,
		Body:    "",
	}
}

// createErrorResponse creates an error response
func (r *RequestHandler) createErrorResponse(reqID, message string, status int) *interfaces.OutgoingResponse {
	return &interfaces.OutgoingResponse{
		ID:      reqID,
		Status:  status,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    message,
	}
}

// isValidHTTPMethod checks if the HTTP method is supported
func (r *RequestHandler) isValidHTTPMethod(method string) bool {
	supportedMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
		"TRACE":   true,
		"CONNECT": true,
	}
	return supportedMethods[method]
}
