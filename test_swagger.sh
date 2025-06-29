#!/bin/bash

# Test script for Swagger initial loading issues
# This script tests the first-time loading of Swagger UI

set -e

echo "🧪 Testing Swagger initial loading..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if tunnel URL is provided
if [ -z "$1" ]; then
    print_error "Usage: $0 <tunnel-url>"
    print_error "Example: $0 https://your-tunnel-id.devpipe.cloud"
    exit 1
fi

TUNNEL_URL="$1"
SWAGGER_URL="${TUNNEL_URL}/docs"

print_status "Testing Swagger at: $SWAGGER_URL"

# Function to test a single request
test_request() {
    local url="$1"
    local description="$2"
    local expected_status="$3"
    
    print_status "Testing: $description"
    
    # Make request and capture response
    response=$(curl -s -w "\n%{http_code}\n%{time_total}\n%{size_download}" \
        -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
        -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8" \
        -H "Accept-Language: en-US,en;q=0.5" \
        -H "Accept-Encoding: gzip, deflate" \
        -H "Connection: keep-alive" \
        -H "Upgrade-Insecure-Requests: 1" \
        "$url" 2>/dev/null)
    
    # Extract status code, time, and size
    status_code=$(echo "$response" | tail -n 2 | head -n 1)
    time_total=$(echo "$response" | tail -n 1)
    size_download=$(echo "$response" | tail -n 3 | head -n 1)
    body=$(echo "$response" | head -n -3)
    
    print_status "  Status: $status_code, Time: ${time_total}s, Size: ${size_download} bytes"
    
    # Check if status code matches expected
    if [ "$status_code" = "$expected_status" ]; then
        print_success "  ✓ Status code $status_code matches expected $expected_status"
    else
        print_error "  ✗ Status code $status_code does not match expected $expected_status"
        return 1
    fi
    
    # Check if response has content
    if [ -n "$body" ]; then
        print_success "  ✓ Response has content"
    else
        print_error "  ✗ Response is empty"
        return 1
    fi
    
    # Check for specific content in Swagger responses
    if [[ "$url" == *"/docs"* ]]; then
        if echo "$body" | grep -q "swagger-ui\|Swagger UI\|OpenAPI"; then
            print_success "  ✓ Contains Swagger UI content"
        else
            print_warning "  ⚠ Response does not contain expected Swagger content"
        fi
    fi
    
    return 0
}

# Function to test asset loading
test_asset() {
    local url="$1"
    local description="$2"
    
    print_status "Testing asset: $description"
    
    response=$(curl -s -w "\n%{http_code}\n%{time_total}\n%{size_download}" \
        -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
        -H "Accept: */*" \
        -H "Accept-Language: en-US,en;q=0.5" \
        -H "Accept-Encoding: gzip, deflate" \
        -H "Referer: $SWAGGER_URL" \
        "$url" 2>/dev/null)
    
    status_code=$(echo "$response" | tail -n 2 | head -n 1)
    time_total=$(echo "$response" | tail -n 1)
    size_download=$(echo "$response" | tail -n 3 | head -n 1)
    
    print_status "  Status: $status_code, Time: ${time_total}s, Size: ${size_download} bytes"
    
    if [ "$status_code" = "200" ] || [ "$status_code" = "304" ]; then
        print_success "  ✓ Asset loaded successfully"
    else
        print_error "  ✗ Asset failed to load (status: $status_code)"
        return 1
    fi
    
    return 0
}

# Test main Swagger page
print_status "=== Testing Main Swagger Page ==="
if test_request "$SWAGGER_URL" "Main Swagger page" "200"; then
    print_success "Main Swagger page loads correctly"
else
    print_error "Main Swagger page failed to load"
    exit 1
fi

# Wait a moment before testing assets
sleep 1

# Test common Swagger assets
print_status "=== Testing Swagger Assets ==="

assets=(
    "/docs/swagger-ui-bundle.js"
    "/docs/swagger-ui-standalone-preset.js"
    "/docs/swagger-ui.css"
    "/docs/swagger-ui-init.js"
    "/docs/favicon-32x32.png"
    "/docs/favicon-16x16.png"
)

failed_assets=0

for asset in "${assets[@]}"; do
    asset_url="${TUNNEL_URL}${asset}"
    if ! test_asset "$asset_url" "$asset"; then
        ((failed_assets++))
    fi
done

# Test with different headers to simulate first-time loading
print_status "=== Testing First-Time Loading Simulation ==="

# Test without cache headers (simulating first visit)
print_status "Testing without cache headers..."
response=$(curl -s -w "\n%{http_code}" \
    -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
    -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8" \
    -H "Cache-Control: no-cache" \
    -H "Pragma: no-cache" \
    "$SWAGGER_URL" 2>/dev/null)

status_code=$(echo "$response" | tail -n 1)
body=$(echo "$response" | head -n -1)

if [ "$status_code" = "200" ] && [ -n "$body" ]; then
    print_success "✓ First-time loading simulation successful"
else
    print_error "✗ First-time loading simulation failed"
fi

# Test with minimal headers (simulating very basic request)
print_status "Testing with minimal headers..."
response=$(curl -s -w "\n%{http_code}" \
    -H "User-Agent: Mozilla/5.0" \
    "$SWAGGER_URL" 2>/dev/null)

status_code=$(echo "$response" | tail -n 1)
body=$(echo "$response" | head -n -1)

if [ "$status_code" = "200" ] && [ -n "$body" ]; then
    print_success "✓ Minimal headers test successful"
else
    print_error "✗ Minimal headers test failed"
fi

# Summary
print_status "=== Test Summary ==="
if [ $failed_assets -eq 0 ]; then
    print_success "All Swagger assets loaded successfully"
else
    print_warning "$failed_assets assets failed to load"
fi

print_status "Swagger initial loading test completed!"
print_status "If you still experience issues with first-time loading, check:"
print_status "1. Network connectivity and latency"
print_status "2. Server response times"
print_status "3. Browser cache settings"
print_status "4. CORS configuration on your local server" 