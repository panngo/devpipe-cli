#!/bin/bash

set -e

echo "🔧 Installing DevPipe..."

ARCH=$(uname -m)
OS=$(uname -s)

# Map architecture
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Map operating system
case "$OS" in
  Linux) PLATFORM="linux" ;;
  Darwin) PLATFORM="darwin" ;;
  *) echo "❌ Unsupported operating system: $OS"; exit 1 ;;
esac

BINARY_URL="https://devpipe.cloud/releases/devpipe-${PLATFORM}-${ARCH}"

echo "➡️  Downloading binary from $BINARY_URL..."
curl -fsSL "$BINARY_URL" -o devpipe
chmod +x devpipe

sudo mv devpipe /usr/local/bin/

echo "✅ DevPipe successfully installed!"
echo "👉 Run: devpipe --help"