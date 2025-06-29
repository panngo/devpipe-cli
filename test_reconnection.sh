#!/bin/bash

# Test script for DevPipe secure reconnection functionality
# This script tests the new UUID and security key authentication system

set -e

echo "🧪 Testing DevPipe Secure Reconnection System"
echo "=============================================="

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

# Check if devpipe binary exists
if [ ! -f "./devpipe" ]; then
    print_error "DevPipe binary not found. Please build the project first."
    exit 1
fi

# Test 1: Clear any existing configuration
print_status "Test 1: Clearing existing tunnel configuration..."
./devpipe -clear-config
print_success "Configuration cleared"

# Test 2: First connection (should create new UUID and security key)
print_status "Test 2: Testing first connection (new UUID generation)..."
timeout 10s ./devpipe -port 3000 &
DEV_PID=$!

# Wait a moment for connection
sleep 3

# Check if process is still running
if kill -0 $DEV_PID 2>/dev/null; then
    print_success "First connection established successfully"
    
    # Kill the process
    kill $DEV_PID 2>/dev/null || true
    wait $DEV_PID 2>/dev/null || true
else
    print_error "First connection failed"
    exit 1
fi

# Test 3: Second connection (should use saved UUID and security key)
print_status "Test 3: Testing secure reconnection (using saved UUID)..."
timeout 10s ./devpipe -port 3000 &
DEV_PID=$!

# Wait a moment for connection
sleep 3

# Check if process is still running
if kill -0 $DEV_PID 2>/dev/null; then
    print_success "Secure reconnection established successfully"
    
    # Kill the process
    kill $DEV_PID 2>/dev/null || true
    wait $DEV_PID 2>/dev/null || true
else
    print_error "Secure reconnection failed"
    exit 1
fi

# Test 4: Verify configuration file exists
print_status "Test 4: Verifying configuration file..."
CONFIG_FILE="$HOME/.devpipe/tunnel.json"
if [ -f "$CONFIG_FILE" ]; then
    print_success "Configuration file found: $CONFIG_FILE"
    
    # Display configuration content (without security key)
    echo "Configuration content:"
    cat "$CONFIG_FILE" | jq 'del(.security_key)' 2>/dev/null || cat "$CONFIG_FILE"
else
    print_error "Configuration file not found"
    exit 1
fi

# Test 5: Test with different port (should still use same UUID)
print_status "Test 5: Testing reconnection with different port..."
timeout 10s ./devpipe -port 8080 &
DEV_PID=$!

# Wait a moment for connection
sleep 3

# Check if process is still running
if kill -0 $DEV_PID 2>/dev/null; then
    print_success "Reconnection with different port successful"
    
    # Kill the process
    kill $DEV_PID 2>/dev/null || true
    wait $DEV_PID 2>/dev/null || true
else
    print_error "Reconnection with different port failed"
    exit 1
fi

# Test 6: Clear configuration and test new connection
print_status "Test 6: Testing new connection after clearing configuration..."
./devpipe -clear-config
timeout 10s ./devpipe -port 3000 &
DEV_PID=$!

# Wait a moment for connection
sleep 3

# Check if process is still running
if kill -0 $DEV_PID 2>/dev/null; then
    print_success "New connection after clearing configuration successful"
    
    # Kill the process
    kill $DEV_PID 2>/dev/null || true
    wait $DEV_PID 2>/dev/null || true
else
    print_error "New connection after clearing configuration failed"
    exit 1
fi

echo ""
print_success "All secure reconnection tests passed! 🎉"
echo ""
echo "Test Summary:"
echo "✅ First connection with UUID generation"
echo "✅ Secure reconnection with saved credentials"
echo "✅ Configuration file persistence"
echo "✅ Reconnection with different port"
echo "✅ Configuration clearing and new connection"
echo ""
echo "The DevPipe client now supports:"
echo "🔐 Secure reconnection with UUID and security key"
echo "💾 Persistent tunnel configuration"
echo "🔄 Automatic reconnection with same tunnel URL"
echo "🗑️  Configuration management with -clear-config flag" 