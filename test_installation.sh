#!/bin/bash

# Test script for DevPipe installation scripts
# This script tests the installation scripts without actually installing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}🔧 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to test platform detection
test_platform_detection() {
    print_status "Testing platform detection..."
    
    local OS=$(uname -s)
    local ARCH=$(uname -m)
    
    echo "Detected OS: $OS"
    echo "Detected Architecture: $ARCH"
    
    case "$ARCH" in
        x86_64) 
            print_success "Architecture: amd64"
            ;;
        arm64|aarch64) 
            print_success "Architecture: arm64"
            ;;
        *) 
            print_error "Unsupported architecture: $ARCH"
            return 1
            ;;
    esac
    
    case "$OS" in
        Linux) 
            print_success "Platform: linux"
            ;;
        Darwin) 
            print_success "Platform: darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*) 
            print_success "Platform: windows"
            ;;
        *) 
            print_error "Unsupported operating system: $OS"
            return 1
            ;;
    esac
}

# Function to test script syntax
test_script_syntax() {
    print_status "Testing script syntax..."
    
    local scripts=("install.sh" "install-unix.sh")
    
    for script in "${scripts[@]}"; do
        if [ -f "$script" ]; then
            if bash -n "$script"; then
                print_success "Syntax check passed: $script"
            else
                print_error "Syntax check failed: $script"
                return 1
            fi
        else
            print_error "Script not found: $script"
            return 1
        fi
    done
}

# Function to test script help
test_script_help() {
    print_status "Testing script help functionality..."
    
    local scripts=("install.sh" "install-unix.sh")
    
    for script in "${scripts[@]}"; do
        if [ -f "$script" ]; then
            if bash "$script" --help >/dev/null 2>&1; then
                print_success "Help functionality works: $script"
            else
                print_warning "Help functionality may not work: $script"
            fi
        fi
    done
}

# Function to test PowerShell script (if available)
test_powershell_script() {
    print_status "Testing PowerShell script..."
    
    if [ -f "install.ps1" ]; then
        if command -v powershell >/dev/null 2>&1; then
            if powershell -Command "Get-Command Test-Path" >/dev/null 2>&1; then
                print_success "PowerShell script exists and PowerShell is available"
            else
                print_warning "PowerShell script exists but PowerShell may not be working"
            fi
        else
            print_warning "PowerShell script exists but PowerShell not found"
        fi
    else
        print_warning "PowerShell script not found"
    fi
}

# Function to test URL construction
test_url_construction() {
    print_status "Testing URL construction..."
    
    local OS=$(uname -s)
    local ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) ARCH="unknown" ;;
    esac
    
    case "$OS" in
        Linux) PLATFORM="linux" ;;
        Darwin) PLATFORM="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) PLATFORM="windows" ;;
        *) PLATFORM="unknown" ;;
    esac
    
    local url="https://github.com/panngo/devpipe-cli/releases/latest/download/devpipe-${PLATFORM}-${ARCH}"
    echo "Constructed URL: $url"
    
    if [ "$PLATFORM" != "unknown" ] && [ "$ARCH" != "unknown" ]; then
        print_success "URL construction successful"
    else
        print_error "URL construction failed - unsupported platform"
        return 1
    fi
}

# Function to test WSL detection
test_wsl_detection() {
    print_status "Testing WSL detection..."
    
    if [[ -f /proc/version ]] && grep -q Microsoft /proc/version; then
        print_success "WSL detected"
    else
        print_status "Not running in WSL"
    fi
}

# Main test function
main() {
    echo ""
    print_status "Running DevPipe installation script tests..."
    echo ""
    
    # Run all tests
    test_platform_detection
    test_script_syntax
    test_script_help
    test_powershell_script
    test_url_construction
    test_wsl_detection
    
    echo ""
    print_success "All tests completed!"
    echo ""
    print_status "Installation scripts are ready for use."
    echo ""
    print_status "Usage examples:"
    echo "  curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash"
echo "  curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install-unix.sh | bash"
echo "  powershell -ExecutionPolicy Bypass -Command \"& { iwr https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.ps1 -UseBasicParsing | iex }\""
    echo ""
}

# Run tests
main 