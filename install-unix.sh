#!/bin/bash

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
    local ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) print_error "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    case "$OS" in
        Linux) PLATFORM="linux" ;;
        Darwin) PLATFORM="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) PLATFORM="windows" ;;
        *) print_error "Unsupported operating system: $OS"; exit 1 ;;
    esac
    
    echo "${PLATFORM}-${ARCH}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to download binary
download_binary() {
    local platform=$1
    local temp_dir=$(mktemp -d)
    local binary_name="devpipe"
    
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="devpipe.exe"
    fi
    
    local binary_url="https://github.com/panngo/devpipe-cli/releases/latest/download/devpipe-${platform}"
    
    print_status "Downloading DevPipe binary for ${platform}..."
    print_status "URL: $binary_url"
    
    if command_exists curl; then
        if curl -fsSL -L "$binary_url" -o "$temp_dir/$binary_name"; then
            print_success "Download completed successfully"
        else
            print_error "Failed to download binary from $binary_url"
            print_error "HTTP Status: $(curl -s -L -o /dev/null -w "%{http_code}" "$binary_url")"
            exit 1
        fi
    elif command_exists wget; then
        if wget -qO "$temp_dir/$binary_name" "$binary_url"; then
            print_success "Download completed successfully"
        else
            print_error "Failed to download binary from $binary_url"
            exit 1
        fi
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    if [ ! -f "$temp_dir/$binary_name" ]; then
        print_error "Binary file not found after download"
        print_error "Temp directory contents:"
        ls -la "$temp_dir/"
        exit 1
    fi
    
    # Check if file is not empty
    if [ ! -s "$temp_dir/$binary_name" ]; then
        print_error "Downloaded file is empty"
        exit 1
    fi
    
    print_success "Binary file verified: $(ls -lh "$temp_dir/$binary_name")"
    
    chmod +x "$temp_dir/$binary_name"
    # Return the path without any colored output
    printf "%s" "$temp_dir/$binary_name"
}

# Function to install binary
install_binary() {
    local binary_path=$1
    local platform=$2
    
    # Determine installation directory
    local install_dir="/usr/local/bin"
    local binary_name="devpipe"
    
    if [[ "$platform" == *"windows"* ]]; then
        install_dir="/usr/local/bin"
        binary_name="devpipe.exe"
    fi
    
    # Check if we have write permissions
    if [ ! -w "$install_dir" ]; then
        print_warning "No write permission to $install_dir, trying with sudo..."
        if command_exists sudo; then
            sudo mv "$binary_path" "$install_dir/$binary_name"
        else
            print_error "No sudo available and no write permission to $install_dir"
            print_error "Please run as root or install sudo"
            exit 1
        fi
    else
        mv "$binary_path" "$install_dir/$binary_name"
    fi
    
    print_success "DevPipe installed to $install_dir/$binary_name"
}

# Function to verify installation
verify_installation() {
    if command_exists devpipe; then
        print_success "DevPipe successfully installed!"
        echo ""
        print_status "Version information:"
        devpipe -help
        echo ""
        print_status "Run 'devpipe --help' to see available commands"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Main installation process
main() {
    echo ""
    print_status "Installing DevPipe..."
    echo ""
    
    # Detect platform
    local platform=$(detect_platform)
    print_status "Detected platform: $platform"
    
    # Download binary (capture only the last line which is the path)
    local binary_path=$(download_binary "$platform" 2>&1 | tail -n1)
    
    # Debug: show the actual path
    print_status "Binary path: '$binary_path'"
    
    # Install binary
    install_binary "$binary_path" "$platform"
    
    # Verify installation
    verify_installation
    
    echo ""
    print_success "Installation completed successfully!"
    print_status "You can now use DevPipe to expose your local applications to the internet."
    echo ""
    print_status "Example usage:"
    echo "  devpipe -port 3000"
    echo "  devpipe -port 8080 -subdomain myapp"
    echo ""
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "DevPipe Installer"
        echo ""
        echo "Usage:"
echo "  curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash"
echo "  curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash -s -- --help"
        echo ""
        echo "Options:"
        echo "  --help, -h    Show this help message"
        echo ""
        echo "Examples:"
echo "  # Install DevPipe"
echo "  curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash"
        echo ""
        echo "  # Install with custom installation directory"
echo "  curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | INSTALL_DIR=/opt/bin bash"
        echo ""
        exit 0
        ;;
    *)
        main
        ;;
esac