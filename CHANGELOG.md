# Changelog

All notable changes to the DevPipe CLI project will be documented in this file.

## [2.0.0] - 2025-06-28

### 🆕 Added

#### Secure Reconnection System
- **UUID Persistence**: Each tunnel now has a unique UUID that persists across reconnections
- **Security Key Authentication**: 32-byte security key for authenticating reconnections
- **Automatic Configuration Management**: UUID and security key are automatically saved and loaded
- **Secure Reconnection**: Clients can reconnect with the same tunnel URL using saved credentials

#### Configuration Management
- **Persistent Storage**: Tunnel configuration saved to `~/.devpipe/tunnel.json`
- **Automatic Loading**: Configuration loaded automatically on startup
- **Manual Configuration Clearing**: New `-clear-config` flag to clear saved configuration
- **Invalid Configuration Handling**: Automatic cleanup of invalid configurations

#### Enhanced UI
- **Security Status Display**: Shows secure reconnection status in banner
- **UUID Display**: Displays current UUID in connection info
- **Security Indicators**: Visual indicators for secure reconnection features

#### New Documentation
- **SECURE_RECONNECTION.md**: Comprehensive documentation of the secure reconnection system
- **Updated README.md**: Complete documentation of new features
- **Test Script**: `test_reconnection.sh` for testing secure reconnection functionality

#### Swagger UI Optimization
- **Enhanced Request Headers**: Automatic addition of proper headers for Swagger requests
- **Improved Response Headers**: Optimized cache control and CORS headers for Swagger assets
- **Content-Type Validation**: Automatic Content-Type correction for JS, CSS, PNG, and ICO files
- **First-Time Loading Fix**: Resolved issue where Swagger UI required page reload on first visit
- **Asset Caching**: Proper cache headers for static assets (1 year) vs dynamic content (no-cache)
- **Detailed Logging**: Enhanced logging for Swagger requests and responses
- **Test Script**: `test_swagger.sh` for testing Swagger loading functionality

#### Transparent Proxy Mode
- **Complete Header Preservation**: All request and response headers are passed through exactly as received
- **No Header Modification**: Removed all header manipulation and optimization logic
- **Transparent Data Flow**: Request and response bodies are passed through without modification
- **Simplified Architecture**: Removed complex header processing and validation logic
- **Universal Compatibility**: Works with any application without requiring specific optimizations
- **Performance Improvement**: Reduced overhead by eliminating header processing

### 🔧 Changed

#### WebSocket Connection Protocol
- **Enhanced Registration**: Registration now supports UUID and security key parameters
- **Improved Error Handling**: Better error messages for authentication failures
- **Thread-Safe Operations**: Enhanced thread safety for WebSocket operations

#### Reconnection Logic
- **Secure Reconnection Priority**: Attempts secure reconnection before falling back to new registration
- **Automatic Fallback**: Falls back to new connection if secure reconnection fails
- **Configuration Validation**: Validates saved configuration before attempting reconnection

#### Client Architecture
- **New Config Package**: Dedicated package for configuration management
- **Enhanced SafeConn**: Added UUID and security key fields to connection struct
- **Improved Error Recovery**: Better error handling and recovery mechanisms

### 🐛 Fixed

- **Configuration Persistence**: Fixed issues with configuration not being saved properly
- **Reconnection Reliability**: Improved reliability of reconnection attempts
- **Error Handling**: Better handling of server errors and network issues

### 📚 Documentation

- **Complete API Documentation**: Documented all new functions and types
- **Usage Examples**: Added comprehensive usage examples
- **Troubleshooting Guide**: Added troubleshooting section for common issues
- **Security Best Practices**: Added security recommendations and best practices

### 🧪 Testing

- **Comprehensive Test Suite**: Added `test_reconnection.sh` for testing secure reconnection
- **Automated Testing**: Updated Makefile with new test targets
- **Test Coverage**: Added tests for all new functionality

## [1.0.0] - Previous Release

### Features
- Basic tunnel functionality
- Automatic reconnection
- Heartbeat system
- HTML/Swagger support
- Next.js optimization

### Technical Details
- Go 1.22+ support
- WebSocket-based communication
- Thread-safe operations
- Comprehensive error handling

---

## Migration Guide

### From v1.0.0 to v2.0.0

#### Automatic Migration
- Existing clients will continue to work without changes
- New secure reconnection features are automatically enabled
- Configuration is automatically created on first run

#### Manual Configuration
- Use `-clear-config` to clear any existing configuration
- Configuration files are stored in `~/.devpipe/tunnel.json`
- No manual migration required

#### Breaking Changes
- None - all changes are backward compatible
- New features are opt-in and don't affect existing functionality

## Security Considerations

### New Security Features
- **UUID Persistence**: Prevents tunnel hijacking
- **Security Key Authentication**: Ensures only authorized clients can reconnect
- **Automatic Cleanup**: Removes inactive tunnels and keys
- **Secure Storage**: Configuration files have restricted permissions

### Best Practices
- Clear configuration periodically for key rotation
- Monitor logs for security-related events
- Don't share configuration files between environments
- Use `-clear-config` when switching environments

## Performance Impact

### Minimal Overhead
- UUID and security key operations are lightweight
- Configuration file operations are infrequent
- No impact on request/response performance
- Automatic cleanup prevents resource leaks

### Benefits
- Improved reliability through secure reconnection
- Better user experience with persistent tunnel URLs
- Enhanced security without performance degradation
- Automatic error recovery and fallback mechanisms 