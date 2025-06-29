.PHONY: help run dev start build docker-start clean test test-reconnection test-all test-http-methods test-third-party test-swagger

# Available commands:
help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the Go client
	go mod tidy && go run main.go

dev: ## Run the Go client (dev mode)
	go mod tidy && go run main.go

start: ## Run the Go client (start mode)
	go mod tidy && go run main.go

build: ## Compile the Go client
	go mod tidy && go build -o devpipe main.go

docker-start: ## Build and run the local Docker image
	docker build -t devpipe-cli .
	docker run -it --rm devpipe-cli

clean: ## Remove the build directory and binary
	rm -rf dist
	rm -f devpipe

test: ## Run all tests
	./test_nextjs.sh
	./test_swagger.sh
	./test_concurrency.sh

test-reconnection: ## Test secure reconnection functionality
	./test_reconnection.sh

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
	./test_nextjs.sh
	./test_swagger.sh
	./test_concurrency.sh
	./test_reconnection.sh
	./test_http_methods.sh
	./test_third_party.sh