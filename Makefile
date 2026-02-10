.PHONY: help build test clean coverage fmt lint install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	go build -ldflags="-s -w" -o color-mcp

test: ## Run tests
	go test -v -race -cover ./...

coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt: ## Format code
	gofmt -w .
	goimports -w .

lint: ## Run linters
	go vet ./...
	staticcheck ./...

install: ## Install to $GOPATH/bin
	go install

clean: ## Clean build artifacts
	rm -f color-mcp
	rm -f coverage.out coverage.html

test-all: fmt lint test ## Run all checks

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	@for GOOS in linux darwin windows; do \
		for GOARCH in amd64 arm64; do \
			NAME=color-mcp-$${GOOS}-$${GOARCH}; \
			if [ "$${GOOS}" = "windows" ]; then \
				NAME=$${NAME}.exe; \
			fi; \
			echo "Building $${NAME}..."; \
			GOOS=$${GOOS} GOARCH=$${GOARCH} CGO_ENABLED=0 \
				go build -ldflags="-s -w" -o dist/$${NAME}; \
		done; \
	done
