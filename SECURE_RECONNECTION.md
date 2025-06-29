# DevPipe Secure Reconnection System

## Overview

The DevPipe client now supports secure reconnection with persistent UUID and security key authentication, allowing you to maintain the same tunnel URL even when the connection drops, with additional authentication for enhanced security.

## 🔐 Security Features

### How It Works

1. **First Connection**: Server generates unique UUID + security key
2. **Reconnection**: Client must provide UUID + valid security key
3. **Validation**: Server verifies the key before authorizing reconnection
4. **Automatic Cleanup**: Inactive tunnels are removed after 1 hour

### Security Key Characteristics

- **Generation**: 32-byte random key (64 hex characters)
- **Storage**: Associated with UUID on server
- **Validation**: Verified on each reconnection attempt
- **Cleanup**: Automatically removed with inactive tunnels

### Security Benefits

1. **Authentication**: Only clients with the correct key can reconnect
2. **Hijacking Prevention**: Prevents other clients from using your UUID
3. **Isolation**: Each tunnel has its own unique key
4. **Expiration**: Keys are removed with inactive tunnels

## Implementation

### Client Implementation

The client automatically handles secure reconnection:

```go
// First connection
registration := map[string]string{
    "action": "register",
    "port":   "3000",
}

// Server responds with UUID and security key
// {
//   "tunnel": "abc123-def456-789-3000",
//   "uuid": "abc123-def456-789", 
//   "key": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234"
// }

// Secure reconnection
registration := map[string]string{
    "action": "register", 
    "port":   "3000",
    "uuid":   "abc123-def456-789",
    "key":    "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234"
}
```

### Configuration Management

The client automatically manages tunnel configuration:

```go
type TunnelConfig struct {
    UUID        string `json:"uuid"`
    SecurityKey string `json:"security_key"`
    TunnelID    string `json:"tunnel_id"`
    Port        string `json:"port"`
}
```

### File Storage

Configuration is stored in `~/.devpipe/tunnel.json`:

```json
{
  "uuid": "abc123-def456-789",
  "security_key": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234",
  "tunnel_id": "abc123-def456-789-3000",
  "port": "3000"
}
```

## Usage

### Command Line Options

```bash
# First connection (generates new UUID)
./devpipe -port 3000

# Reconnection (uses saved UUID)
./devpipe -port 3000

# Clear saved configuration
./devpipe -clear-config

# Clear configuration and specify port
./devpipe -clear-config -port 8080
```

### Automatic Behavior

1. **First Run**: Creates new tunnel with UUID and security key
2. **Subsequent Runs**: Automatically uses saved UUID and key
3. **Invalid Configuration**: Clears invalid config and creates new connection
4. **Manual Clear**: Use `-clear-config` to force new connection

## Error Handling

### Server Error Responses

If the security key is invalid or missing:

```json
{
  "error": "Invalid security key"
}
```

```json
{
  "error": "Security key required for reconnection"
}
```

### Client Error Handling

The client handles various error scenarios:

1. **Invalid Security Key**: Clears configuration and creates new connection
2. **Missing Configuration**: Creates new connection automatically
3. **Server Errors**: Logs error and attempts fallback
4. **Network Errors**: Retries with exponential backoff

## Logs and Debugging

### Security Logs

The client provides detailed security logs:

```
🔐 Attempting secure reconnection with UUID: abc123-def456-789
💾 Tunnel configuration saved for secure reconnection
🗑️  Cleared invalid tunnel configuration
🔑 UUID: abc123-def456-789
```

### Debugging Commands

```bash
# View security logs
./devpipe -port 3000 2>&1 | grep -E "(🔐|🔑|💾|🗑️)"

# View configuration file
cat ~/.devpipe/tunnel.json

# Clear configuration
./devpipe -clear-config
```

## Testing

### Test Script

Run the comprehensive test script:

```bash
./test_reconnection.sh
```

This script tests:
- ✅ First connection with UUID generation
- ✅ Secure reconnection with saved credentials
- ✅ Configuration file persistence
- ✅ Reconnection with different port
- ✅ Configuration clearing and new connection

### Manual Testing

1. **First Connection**:
   ```bash
   ./devpipe -port 3000
   # Note the UUID in logs
   ```

2. **Reconnection**:
   ```bash
   ./devpipe -port 3000
   # Should show "🔐 Attempting secure reconnection"
   ```

3. **Clear Configuration**:
   ```bash
   ./devpipe -clear-config
   # Should show "🗑️  Tunnel configuration cleared"
   ```

## Best Practices

### Security

1. **Secure Storage**: Configuration file has 600 permissions
2. **Key Rotation**: Clear configuration periodically for new keys
3. **Environment Isolation**: Each environment should have separate configuration
4. **Monitoring**: Monitor logs for security-related events

### Reliability

1. **Automatic Reconnection**: Client handles reconnection automatically
2. **Fallback Strategy**: Falls back to new connection if reconnection fails
3. **Error Recovery**: Clears invalid configuration automatically
4. **Logging**: Comprehensive logging for troubleshooting

### Configuration Management

1. **Backup**: Backup configuration files if needed
2. **Version Control**: Don't commit configuration files to version control
3. **Cleanup**: Use `-clear-config` when switching environments
4. **Monitoring**: Monitor configuration file changes

## Troubleshooting

### Common Issues

1. **"Invalid security key" error**:
   - Clear configuration: `./devpipe -clear-config`
   - Check if server was restarted
   - Verify tunnel timeout (1 hour)

2. **Configuration not saved**:
   - Check file permissions on `~/.devpipe/`
   - Verify disk space
   - Check for write errors in logs

3. **Reconnection not working**:
   - Clear configuration and try again
   - Check network connectivity
   - Verify server is running

4. **UUID mismatch**:
   - This is normal when server assigns new UUID
   - Check logs for "⚠️  New tunnel ID assigned"
   - Configuration will be updated automatically

### Debugging Steps

1. **Check Configuration**:
   ```bash
   cat ~/.devpipe/tunnel.json
   ```

2. **View Logs**:
   ```bash
   ./devpipe -port 3000 2>&1 | grep -E "(🔐|🔑|💾|🗑️|⚠️|❌)"
   ```

3. **Clear and Retry**:
   ```bash
   ./devpipe -clear-config
   ./devpipe -port 3000
   ```

4. **Check File Permissions**:
   ```bash
   ls -la ~/.devpipe/
   ```

## Compatibility

- ✅ Existing clients continue working (new connection without key)
- ✅ UUID is optional - if not sent, a new one will be generated
- ✅ Thread-safe system for multiple simultaneous connections
- ✅ Key required only for reconnections

## Future Enhancements

1. **Key Rotation**: Automatic key rotation for enhanced security
2. **Multiple Tunnels**: Support for multiple tunnels per client
3. **Configuration Encryption**: Encrypt configuration file
4. **Audit Logging**: Enhanced logging for security events
5. **Rate Limiting**: Prevent brute force attacks on security keys 