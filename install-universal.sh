#!/bin/bash

# DevPipe Universal Installer
# This script automatically detects the platform and uses the appropriate installation method
# Usage: curl -fsSL https://devpipe.cloud/install.sh | bash

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

# Function to detect platform
detect_platform() {
    local OS=$(uname -s)
    
    case "$OS" in
        Linux) echo "linux" ;;
        Darwin) echo "macos" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

# Function to check if running in WSL
is_wsl() {
    if [[ -f /proc/version ]] && grep -q Microsoft /proc/version; then
        return 0
    else
        return 1
    fi
}

# Function to install on Linux/macOS
install_unix() {
    print_status "Installing DevPipe on Unix-like system..."
    
    # Download and execute the Unix installer
    curl -fsSL https://devpipe.cloud/install-unix.sh | bash
}

# Function to install on Windows
install_windows() {
    print_status "Installing DevPipe on Windows..."
    
    if command -v powershell >/dev/null 2>&1; then
        print_status "Using PowerShell installer..."
        powershell -ExecutionPolicy Bypass -Command "& { iwr https://devpipe.cloud/install.ps1 -UseBasicParsing | iex }"
    else
        print_error "PowerShell not found. Please install PowerShell and try again."
        exit 1
    fi
}

# Function to install on WSL
install_wsl() {
    print_status "Installing DevPipe on WSL (Windows Subsystem for Linux)..."
    print_warning "WSL detected. Using Unix installer..."
    
    # Use the Unix installer for WSL
    install_unix
}

# Main installation process
main() {
    echo ""
    print_status "DevPipe Universal Installer"
    echo ""
    
    # Detect platform
    local platform=$(detect_platform)
    print_status "Detected platform: $platform"
    
    case "$platform" in
        "linux")
            if is_wsl; then
                install_wsl
            else
                install_unix
            fi
            ;;
        "macos")
            install_unix
            ;;
        "windows")
            install_windows
            ;;
        "unknown")
            print_error "Unsupported platform: $(uname -s)"
            print_error "Please install manually or contact support"
            exit 1
            ;;
    esac
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "DevPipe Universal Installer"
        echo ""
        echo "This installer automatically detects your platform and uses the appropriate installation method."
        echo ""
        echo "Usage:"
        echo "  curl -fsSL https://devpipe.cloud/install.sh | bash"
        echo "  curl -fsSL https://devpipe.cloud/install.sh | bash -s -- --help"
        echo ""
        echo "Supported platforms:"
        echo "  - Linux (including WSL)"
        echo "  - macOS"
        echo "  - Windows (PowerShell required)"
        echo ""
        echo "Examples:"
        echo "  # Install DevPipe (automatic platform detection)"
        echo "  curl -fsSL https://devpipe.cloud/install.sh | bash"
        echo ""
        echo "  # Install with help"
        echo "  curl -fsSL https://devpipe.cloud/install.sh | bash -s -- --help"
        echo ""
        echo "Manual installation:"
        echo "  # Linux/macOS"
        echo "  curl -fsSL https://devpipe.cloud/install-unix.sh | bash"
        echo ""
        echo "  # Windows"
        echo "  powershell -ExecutionPolicy Bypass -Command \"& { iwr https://devpipe.cloud/install.ps1 -UseBasicParsing | iex }\""
        echo ""
        exit 0
        ;;
    *)
        main
        ;;
esac 