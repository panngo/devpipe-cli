# DevPipe Windows Installer
# Usage: powershell -ExecutionPolicy Bypass -Command "& { iwr https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.ps1 -UseBasicParsing | iex }"

param(
    [switch]$Help
)

# Function to write colored output
function Write-Status {
    param([string]$Message)
    Write-Host "🔧 $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor Red
}

# Function to detect platform
function Get-Platform {
    $arch = $env:PROCESSOR_ARCHITECTURE
    $os = $env:OS
    
    switch ($arch) {
        "AMD64" { $arch = "amd64" }
        "ARM64" { $arch = "arm64" }
        default { 
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
    
    if ($os -like "*Windows*") {
        $platform = "windows"
    } else {
        Write-Error "Unsupported operating system"
        exit 1
    }
    
    return "$platform-$arch"
}

# Function to download binary
function Get-DevPipeBinary {
    param([string]$Platform)
    
    $tempDir = [System.IO.Path]::GetTempPath()
    $binaryName = "devpipe.exe"
    $binaryUrl = "https://github.com/panngo/devpipe-cli/releases/latest/download/devpipe-$Platform"
    $binaryPath = Join-Path $tempDir $binaryName
    
    Write-Status "Downloading DevPipe binary for $Platform..."
    Write-Status "URL: $binaryUrl"
    
    try {
        Invoke-WebRequest -Uri $binaryUrl -OutFile $binaryPath -UseBasicParsing
        Write-Success "Binary downloaded successfully"
        return $binaryPath
    }
    catch {
        Write-Error "Failed to download binary: $($_.Exception.Message)"
        exit 1
    }
}

# Function to install binary
function Install-DevPipeBinary {
    param([string]$BinaryPath)
    
    # Determine installation directory
    $installDir = "$env:ProgramFiles\DevPipe"
    $binaryName = "devpipe.exe"
    $installPath = Join-Path $installDir $binaryName
    
    Write-Status "Installing DevPipe to $installPath..."
    
    try {
        # Create directory if it doesn't exist
        if (!(Test-Path $installDir)) {
            New-Item -ItemType Directory -Path $installDir -Force | Out-Null
        }
        
        # Copy binary
        Copy-Item -Path $BinaryPath -Destination $installPath -Force
        
        # Add to PATH if not already there
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($currentPath -notlike "*$installDir*") {
            $newPath = "$currentPath;$installDir"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
            Write-Warning "Added DevPipe to PATH. You may need to restart your terminal."
        }
        
        Write-Success "DevPipe installed to $installPath"
        return $installPath
    }
    catch {
        Write-Error "Failed to install binary: $($_.Exception.Message)"
        exit 1
    }
}

# Function to verify installation
function Test-DevPipeInstallation {
    try {
        $version = & "devpipe.exe" --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "DevPipe successfully installed!"
            Write-Host ""
            Write-Status "Version information:"
            Write-Host $version
            Write-Host ""
            Write-Status "Run 'devpipe --help' to see available commands"
            return $true
        } else {
            Write-Error "Installation verification failed"
            return $false
        }
    }
    catch {
        Write-Error "Installation verification failed: $($_.Exception.Message)"
        return $false
    }
}

# Main installation process
function Install-DevPipe {
    Write-Host ""
    Write-Status "Installing DevPipe..."
    Write-Host ""
    
    # Detect platform
    $platform = Get-Platform
    Write-Status "Detected platform: $platform"
    
    # Download binary
    $binaryPath = Get-DevPipeBinary -Platform $platform
    
    # Install binary
    $installPath = Install-DevPipeBinary -BinaryPath $binaryPath
    
    # Verify installation
    $success = Test-DevPipeInstallation
    
    if ($success) {
        Write-Host ""
        Write-Success "Installation completed successfully!"
        Write-Status "You can now use DevPipe to expose your local applications to the internet."
        Write-Host ""
        Write-Status "Example usage:"
        Write-Host "  devpipe -port 3000"
        Write-Host "  devpipe -port 8080 -subdomain myapp"
        Write-Host ""
    } else {
        Write-Error "Installation failed"
        exit 1
    }
}

# Handle help parameter
if ($Help) {
    Write-Host "DevPipe Windows Installer"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  powershell -ExecutionPolicy Bypass -Command `"& { iwr https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.ps1 -UseBasicParsing | iex }`""
    Write-Host "  powershell -ExecutionPolicy Bypass -File install.ps1 -Help"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Help    Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  # Install DevPipe"
    Write-Host "  powershell -ExecutionPolicy Bypass -Command `"& { iwr https://raw.githubusercontent.com/panngo/devpipe-cli/main/install.ps1 -UseBasicParsing | iex }`""
    Write-Host ""
    Write-Host "  # Install with help"
    Write-Host "  powershell -ExecutionPolicy Bypass -File install.ps1 -Help"
    Write-Host ""
    exit 0
}

# Run installation
Install-DevPipe 