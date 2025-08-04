package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger interface following Interface Segregation Principle
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	SetLevel(level LogLevel)
	SetOutput(writer io.Writer)
}

// ConsoleLogger implements Logger for console output
type ConsoleLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger(level LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Debug logs debug messages
func (l *ConsoleLogger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, args...)
	}
}

// Info logs info messages
func (l *ConsoleLogger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, args...)
	}
}

// Warn logs warning messages
func (l *ConsoleLogger) Warn(format string, args ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] "+format, args...)
	}
}

// Error logs error messages
func (l *ConsoleLogger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] "+format, args...)
	}
}

// SetLevel sets the logging level
func (l *ConsoleLogger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer
func (l *ConsoleLogger) SetOutput(writer io.Writer) {
	l.logger.SetOutput(writer)
}

// SilentLogger implements Logger for silent operation
type SilentLogger struct {
	level LogLevel
}

// NewSilentLogger creates a new silent logger
func NewSilentLogger() *SilentLogger {
	return &SilentLogger{
		level: ERROR, // Only show errors
	}
}

// Debug logs debug messages (silent)
func (l *SilentLogger) Debug(format string, args ...interface{}) {
	// Silent
}

// Info logs info messages (silent)
func (l *SilentLogger) Info(format string, args ...interface{}) {
	// Silent
}

// Warn logs warning messages (silent)
func (l *SilentLogger) Warn(format string, args ...interface{}) {
	// Silent
}

// Error logs error messages
func (l *SilentLogger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", args...)
	}
}

// SetLevel sets the logging level
func (l *SilentLogger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer (no-op for silent logger)
func (l *SilentLogger) SetOutput(writer io.Writer) {
	// No-op for silent logger
}

// HTTPRequestLogger logs HTTP requests in a structured way
type HTTPRequestLogger struct {
	logger Logger
}

// NewHTTPRequestLogger creates a new HTTP request logger
func NewHTTPRequestLogger(logger Logger) *HTTPRequestLogger {
	return &HTTPRequestLogger{
		logger: logger,
	}
}

// LogRequest logs an HTTP request
func (h *HTTPRequestLogger) LogRequest(method, path, status string, duration time.Duration) {
	h.logger.Info("🌐 HTTP %s %s - %s (%v)", method, path, status, duration)
}

// LogTunnelEvent logs tunnel-related events
func (h *HTTPRequestLogger) LogTunnelEvent(event, tunnelID string) {
	h.logger.Info("🔗 %s: %s", event, tunnelID)
}

// LogReconnectionEvent logs reconnection events
func (h *HTTPRequestLogger) LogReconnectionEvent(event, url string) {
	h.logger.Info("🔄 %s: %s", event, url)
}
