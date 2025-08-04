# DevPipe Installation Guide

This document explains the different installation methods available for DevPipe CLI.

## 🚀 Quick Installation

### Universal Installer (Recommended)

The universal installer automatically detects your platform and uses the appropriate installation method:

```bash
curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash
```

This works on:
- ✅ Linux
- ✅ macOS  
- ✅ Windows (PowerShell required)
- ✅ WSL (Windows Subsystem for Linux)

## 📋 Platform-Specific Installation

### Linux/macOS

For Unix-like systems (Linux, macOS, WSL):

```bash
curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install-unix.sh | bash
```

**Features:**
- Automatic platform detection (x86_64, arm64)
- Downloads from `https://github.com/panngo/devpipe-cli/releases/latest/download/`
- Installs to `/usr/local/bin/`
- Handles permissions with sudo if needed
- Supports both curl and wget

### Windows

For Windows systems (PowerShell required):

```powershell
powershell -ExecutionPolicy Bypass -Command "& { iwr https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.ps1 -UseBasicParsing | iex }"
```

**Features:**
- Installs to `C:\Program Files\DevPipe\`
- Automatically adds to PATH
- Supports x64 and ARM64 architectures
- Handles Windows-specific paths and permissions

## 🔧 Manual Installation

If you prefer to install manually:

### 1. Download Binary

Visit the releases page and download the appropriate binary for your platform:
- `devpipe-linux-amd64` - Linux x64
- `devpipe-linux-arm64` - Linux ARM64
- `devpipe-darwin-amd64` - macOS Intel
- `devpipe-darwin-arm64` - macOS Apple Silicon
- `devpipe-windows-amd64.exe` - Windows x64
- `devpipe-windows-arm64.exe` - Windows ARM64

### 2. Make Executable (Unix-like systems)

```bash
chmod +x devpipe
```

### 3. Install to PATH

**Linux/macOS:**
```bash
sudo mv devpipe /usr/local/bin/
```

**Windows:**
```powershell
# Create directory
New-Item -ItemType Directory -Path "$env:ProgramFiles\DevPipe" -Force

# Move binary
Move-Item devpipe.exe "$env:ProgramFiles\DevPipe\"

# Add to PATH
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$newPath = "$currentPath;$env:ProgramFiles\DevPipe"
[Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
```

## 🐳 Docker Installation

You can also run DevPipe in a Docker container:

```bash
docker run --rm -p 3000:3000 devpipe/devpipe-server
```

## 🔍 Verification

After installation, verify that DevPipe is working:

```bash
devpipe --version
devpipe --help
```

## 🛠 Troubleshooting

### Permission Denied

**Linux/macOS:**
```bash
# Check if binary is executable
ls -la /usr/local/bin/devpipe

# Fix permissions if needed
sudo chmod +x /usr/local/bin/devpipe
```

**Windows:**
```powershell
# Check if binary exists
Test-Path "$env:ProgramFiles\DevPipe\devpipe.exe"

# Check PATH
$env:PATH -split ';' | Where-Object { $_ -like "*DevPipe*" }
```

### Binary Not Found

**Linux/macOS:**
```bash
# Check if binary is in PATH
which devpipe

# Check common locations
ls -la /usr/local/bin/devpipe
ls -la /usr/bin/devpipe
```

**Windows:**
```powershell
# Check if binary is in PATH
Get-Command devpipe.exe

# Check installation directory
Test-Path "$env:ProgramFiles\DevPipe\devpipe.exe"
```

### Network Issues

If you're behind a corporate firewall or proxy:

```bash
# Use wget instead of curl
wget -qO- https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash

# Or download manually and run locally
curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh -o install.sh
bash install.sh
```

## 📝 Installation Scripts

The installation system consists of several scripts:

- `install.sh` - Universal installer (main entry point)
- `install-unix.sh` - Unix-like systems installer
- `install.ps1` - Windows PowerShell installer

### Script Features

**Universal Installer (`install.sh`):**
- Platform detection
- Automatic routing to appropriate installer
- WSL detection and handling
- Help and documentation

**Unix Installer (`install-unix.sh`):**
- Architecture detection (x86_64, arm64)
- Multiple download methods (curl, wget)
- Permission handling
- Installation verification

**Windows Installer (`install.ps1`):**
- PowerShell-based installation
- Windows-specific paths
- PATH environment variable management
- Installation verification

## 🔄 Updating

To update DevPipe, simply run the installation script again:

```bash
curl -fsSL https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.sh | bash
```

The installer will download and install the latest version.

## 🗑 Uninstalling

**Linux/macOS:**
```bash
sudo rm /usr/local/bin/devpipe
```

**Windows:**
```powershell
# Remove binary
Remove-Item "$env:ProgramFiles\DevPipe\devpipe.exe" -Force

# Remove directory if empty
Remove-Item "$env:ProgramFiles\DevPipe" -Force

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$newPath = ($currentPath -split ';' | Where-Object { $_ -notlike "*DevPipe*" }) -join ';'
[Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
```

## 📞 Support

If you encounter issues with installation:

1. Check the troubleshooting section above
2. Verify your platform is supported
3. Check network connectivity
4. Review the installation logs
5. Open an issue on GitHub with:
   - Your platform and architecture
   - Installation command used
   - Error messages
   - Installation logs 