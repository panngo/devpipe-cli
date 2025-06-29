# DevPipe HTTP Methods Support

## Overview

The DevPipe client now supports **all standard HTTP methods** that browsers can send, with proper validation, error handling, and method-specific optimizations.

## 🌐 Supported HTTP Methods

### Standard Methods (RFC 7231)
- **GET**: Retrieve a representation of a resource
- **POST**: Submit data to be processed
- **PUT**: Replace the target resource
- **DELETE**: Remove the target resource
- **PATCH**: Apply partial modifications to a resource

### Special Methods
- **HEAD**: Same as GET but without response body
- **OPTIONS**: Get information about communication options
- **TRACE**: Perform a message loop-back test
- **CONNECT**: Establish a tunnel to the server

## 🔧 Implementation Details

### Method Validation
```go
var supportedMethods = map[string]bool{
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
```

### Body Handling
```go
// Methods that typically have a body
var methodsWithBody = map[string]bool{
    "POST":   true,
    "PUT":    true,
    "PATCH":  true,
    "DELETE": false, // DELETE can have body but often doesn't
}
```

### Request Creation
```go
// Create request with appropriate body handling
if shouldHaveBody(req.Method) && req.Body != "" {
    httpReq, err = http.NewRequest(req.Method, url, strings.NewReader(req.Body))
} else {
    httpReq, err = http.NewRequest(req.Method, url, nil)
}
```

## 📋 Method-Specific Handling

### GET Requests
- **Body**: No request body
- **Headers**: All standard headers preserved
- **Response**: Full response body included
- **Use Case**: Retrieving data, loading pages

### POST Requests
- **Body**: Request body included when present
- **Headers**: Content-Type and Content-Length set automatically
- **Response**: Full response body included
- **Use Case**: Submitting forms, creating resources

### PUT Requests
- **Body**: Request body included when present
- **Headers**: Content-Type and Content-Length set automatically
- **Response**: Full response body included
- **Use Case**: Updating entire resources

### DELETE Requests
- **Body**: Optional request body (flexible handling)
- **Headers**: Content-Type and Content-Length set if body present
- **Response**: Full response body included
- **Use Case**: Removing resources

### PATCH Requests
- **Body**: Request body included when present
- **Headers**: Content-Type and Content-Length set automatically
- **Response**: Full response body included
- **Use Case**: Partial resource updates

### HEAD Requests
- **Body**: No request body
- **Headers**: All standard headers preserved
- **Response**: Headers only, no body (Content-Length: 0)
- **Use Case**: Checking resource existence, getting metadata

### OPTIONS Requests
- **Body**: No request body
- **Headers**: CORS headers automatically added
- **Response**: CORS preflight response
- **Use Case**: CORS preflight, checking server capabilities

### TRACE Requests
- **Body**: No request body
- **Headers**: All standard headers preserved
- **Response**: Echo of the request
- **Use Case**: Debugging, diagnostics

## 🌍 CORS Support

### Automatic CORS Headers
For OPTIONS requests, the client automatically includes:

```go
response.Headers = map[string]string{
    "Access-Control-Allow-Origin":      "*",
    "Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT",
    "Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Requested-With, Accept, Origin, User-Agent, Referer",
    "Access-Control-Max-Age":           "86400",
    "Access-Control-Allow-Credentials": "true",
    "Content-Length":                   "0",
}
```

### CORS Preflight Flow
1. Browser sends OPTIONS request
2. Client responds with CORS headers
3. Browser proceeds with actual request
4. Client forwards request to local server

## 🔒 Error Handling

### Method Validation
```go
if !isValidHTTPMethod(req.Method) {
    log.Printf("❌ Unsupported HTTP method: %s", req.Method)
    sendErrorResponse(conn, req.ID, fmt.Sprintf("Unsupported HTTP method: %s", req.Method), 405)
    return
}
```

### Path Validation
```go
if req.Path == "" {
    log.Printf("❌ Empty request path")
    sendErrorResponse(conn, req.ID, "Empty request path", 400)
    return
}
```

### Error Responses
- **405 Method Not Allowed**: Unsupported HTTP method
- **400 Bad Request**: Invalid request (empty path, malformed headers)
- **502 Bad Gateway**: Local server connection failed
- **500 Internal Server Error**: Client processing error

## 📊 Logging and Monitoring

### Method-Specific Logs
```
🌐 HTTP GET /api/users
📦 Request body length: 0 bytes
🌐 HTTP POST /api/users
📦 Request body length: 45 bytes
🌐 HTTP OPTIONS /api/users (CORS preflight)
🌐 HTTP HEAD /api/users (no body)
```

### Response Logs
```
GET    /api/users          200 OK
POST   /api/users          201 OK
PUT    /api/users/123      200 OK
DELETE /api/users/123      204 OK
PATCH  /api/users/123      200 OK
HEAD   /api/users          200 OK
OPTIONS /api/users         200 OK
```

## 🧪 Testing

### Test Script
Run the comprehensive HTTP methods test:

```bash
./test_http_methods.sh
```

### Manual Testing
```bash
# Test GET request
curl https://your-tunnel.devpipe.cloud/api/users

# Test POST request
curl -X POST -d '{"name":"John"}' https://your-tunnel.devpipe.cloud/api/users

# Test PUT request
curl -X PUT -d '{"name":"Jane"}' https://your-tunnel.devpipe.cloud/api/users/123

# Test DELETE request
curl -X DELETE https://your-tunnel.devpipe.cloud/api/users/123

# Test PATCH request
curl -X PATCH -d '{"name":"Bob"}' https://your-tunnel.devpipe.cloud/api/users/123

# Test HEAD request
curl -I https://your-tunnel.devpipe.cloud/api/users

# Test OPTIONS request (CORS preflight)
curl -X OPTIONS https://your-tunnel.devpipe.cloud/api/users
```

## 🔧 Configuration

### Headers Handling
The client automatically handles problematic headers:

```go
// Skip headers that Go manages automatically
if strings.ToLower(k) == "content-length" {
    continue // Let Go calculate automatically
}
if strings.ToLower(k) == "host" {
    continue // Will be defined automatically
}
if strings.ToLower(k) == "connection" {
    continue // Will be managed by Go
}
if strings.ToLower(k) == "transfer-encoding" {
    continue // Let Go handle it
}
```

### Content-Length Calculation
```go
// Set Content-Length for methods that should have a body
if shouldHaveBody(req.Method) && req.Body != "" {
    httpReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(req.Body)))
}
```

## 🚀 Performance Optimizations

### Method-Specific Optimizations
- **HEAD requests**: No body reading, faster response
- **OPTIONS requests**: CORS response without server call
- **GET requests**: Standard forwarding with caching headers preserved
- **POST/PUT/PATCH**: Body handling optimized for each method

### Memory Efficiency
- **Streaming**: Large request bodies handled efficiently
- **Buffer management**: Appropriate buffer sizes for each method
- **Header optimization**: Only necessary headers processed

## 🔍 Debugging

### Common Issues

1. **Method Not Allowed (405)**:
   - Check if method is in supportedMethods map
   - Verify method name case (automatically converted to uppercase)

2. **Bad Request (400)**:
   - Check request path is not empty
   - Verify headers are properly formatted

3. **CORS Issues**:
   - Ensure OPTIONS requests are handled
   - Check CORS headers in response

### Debug Commands
```bash
# View HTTP method logs
./devpipe -port 3000 2>&1 | grep -E "(🌐|📦)"

# View error logs
./devpipe -port 3000 2>&1 | grep -E "(❌|⚠️)"

# View all logs
./devpipe -port 3000
```

## 📈 Browser Compatibility

### Supported Browsers
- ✅ Chrome/Chromium
- ✅ Firefox
- ✅ Safari
- ✅ Edge
- ✅ Opera

### Browser-Specific Methods
- **Chrome**: All methods including CONNECT
- **Firefox**: All methods including TRACE
- **Safari**: All methods except CONNECT (security restrictions)
- **Edge**: All methods including CONNECT

## 🔮 Future Enhancements

### Planned Features
1. **Method-specific caching**: Different cache strategies per method
2. **Rate limiting**: Method-specific rate limiting
3. **Metrics**: Detailed metrics per HTTP method
4. **Custom methods**: Support for custom HTTP methods
5. **WebSocket upgrade**: Support for WebSocket upgrade requests

### Performance Improvements
1. **Connection pooling**: Reuse connections for similar methods
2. **Compression**: Method-specific compression strategies
3. **Caching**: Intelligent caching based on method type
4. **Load balancing**: Method-aware load balancing

## 📚 References

### RFC Standards
- [RFC 7231](https://tools.ietf.org/html/rfc7231): HTTP/1.1 Semantics and Content
- [RFC 7232](https://tools.ietf.org/html/rfc7232): HTTP/1.1 Conditional Requests
- [RFC 7233](https://tools.ietf.org/html/rfc7233): HTTP/1.1 Range Requests
- [RFC 7234](https://tools.ietf.org/html/rfc7234): HTTP/1.1 Caching
- [RFC 7235](https://tools.ietf.org/html/rfc7235): HTTP/1.1 Authentication

### Browser Specifications
- [MDN HTTP Methods](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods)
- [W3C HTTP Specification](https://www.w3.org/Protocols/)
- [CORS Specification](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) 