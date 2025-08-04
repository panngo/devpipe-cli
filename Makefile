.PHONY: help run dev dev-local dev-debug start build build-all build-linux build-darwin build-windows docker-start clean test test-reconnection test-all test-http-methods test-third-party test-swagger install uninstall check fmt

# Available commands:
help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the Go client
	go mod tidy && go run main.go

dev: ## Run the Go client (dev mode)
	go mod tidy && go run main.go

dev-local: ## Run the Go client with local server (dev mode)
	go mod tidy && go run main.go -port 3000

start: ## Run the Go client (start mode)
	go mod tidy && go run main.go

dev-debug: ## Run the Go client with debug mode
	go mod tidy && go run main.go -debug

build: ## Compile the Go client
	go mod tidy && go build -o devpipe main.go
	@echo "✅ Binary built: ./devpipe"

build-all: ## Build binaries for all platforms
	@echo "🔨 Building binaries for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -o dist/devpipe-linux-amd64 main.go
	@GOOS=linux GOARCH=arm64 go build -o dist/devpipe-linux-arm64 main.go
	@GOOS=darwin GOARCH=amd64 go build -o dist/devpipe-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 go build -o dist/devpipe-darwin-arm64 main.go
	@GOOS=windows GOARCH=amd64 go build -o dist/devpipe-windows-amd64.exe main.go
	@GOOS=windows GOARCH=arm64 go build -o dist/devpipe-windows-arm64.exe main.go
	@chmod +x dist/devpipe-linux-amd64 dist/devpipe-linux-arm64 dist/devpipe-darwin-amd64 dist/devpipe-darwin-arm64
	@echo "✅ All binaries built successfully in dist/ directory"
	@echo "📦 Binaries:"
	@ls -lh dist/

build-linux: ## Build Linux binaries
	@echo "🔨 Building Linux binaries..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -o dist/devpipe-linux-amd64 main.go
	@GOOS=linux GOARCH=arm64 go build -o dist/devpipe-linux-arm64 main.go
	@chmod +x dist/devpipe-linux-amd64 dist/devpipe-linux-arm64
	@echo "✅ Linux binaries built successfully"
	@ls -lh dist/devpipe-linux-*

build-darwin: ## Build macOS binaries
	@echo "🔨 Building macOS binaries..."
	@mkdir -p dist
	@GOOS=darwin GOARCH=amd64 go build -o dist/devpipe-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 go build -o dist/devpipe-darwin-arm64 main.go
	@chmod +x dist/devpipe-darwin-amd64 dist/devpipe-darwin-arm64
	@echo "✅ macOS binaries built successfully"
	@ls -lh dist/devpipe-darwin-*

build-windows: ## Build Windows binaries
	@echo "🔨 Building Windows binaries..."
	@mkdir -p dist
	@GOOS=windows GOARCH=amd64 go build -o dist/devpipe-windows-amd64.exe main.go
	@GOOS=windows GOARCH=arm64 go build -o dist/devpipe-windows-arm64.exe main.go
	@echo "✅ Windows binaries built successfully"
	@ls -lh dist/devpipe-windows-*

docker-start: ## Build and run the local Docker image
	docker build -t devpipe-cli .
	docker run -it --rm devpipe-cli

clean: ## Remove the build directory and binary
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf dist
	@rm -f devpipe
	@echo "✅ Clean completed"

install: build ## Install devpipe to /usr/local/bin
	@echo "📦 Installing devpipe to /usr/local/bin..."
	@sudo cp devpipe /usr/local/bin/
	@echo "✅ DevPipe installed successfully"
	@echo "🚀 Usage: devpipe -help"

uninstall: ## Remove devpipe from /usr/local/bin
	@echo "🗑️ Uninstalling devpipe from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/devpipe
	@echo "✅ DevPipe uninstalled successfully"

test: ## Run all tests
	@echo "🧪 Running all tests..."
	@./test_nextjs.sh
	@./test_swagger.sh
	@./test_concurrency.sh
	@echo "✅ All tests completed"

test-reconnection: ## Test secure reconnection functionality
	@echo "🔌 Testing secure reconnection..."
	@./test_reconnection.sh
	@echo "✅ Reconnection test completed"

test-http-methods: build
	@echo "🧪 Testing HTTP methods support..."
	@chmod +x test_http_methods.sh
	@./test_http_methods.sh

test-third-party: build
	@echo "🧪 Testing third-party request handling..."
	@chmod +x test_third_party.sh
	@./test_third_party.sh

test-swagger: build
	@echo "🧪 Testing Swagger initial loading..."
	@chmod +x test_swagger.sh
	@echo "Usage: make test-swagger TUNNEL_URL=https://your-tunnel-id.devpipe.cloud"
	@echo "Example: make test-swagger TUNNEL_URL=https://11c926be-39d9-484b-b409-659248402687-3003.devpipe.cloud"

test-all: ## Run all tests including secure reconnection, HTTP methods, and third-party
	@echo "🧪 Running comprehensive test suite..."
	@./test_nextjs.sh
	@./test_swagger.sh
	@./test_concurrency.sh
	@./test_reconnection.sh
	@./test_http_methods.sh
	@./test_third_party.sh
	@echo "✅ All tests completed successfully"

check: ## Check code quality and dependencies
	@echo "🔍 Checking code quality..."
	@go mod tidy
	@go mod verify
	@go vet ./...
	@echo "✅ Code quality check completed"

fmt: ## Format Go code
	@echo "🎨 Formatting Go code..."
	@go fmt ./...
	@echo "✅ Code formatting completed"