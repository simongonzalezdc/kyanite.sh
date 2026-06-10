.PHONY: test build install clean lint fmt run help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

build: ## Build the binary
	go build -o bin/prism ./cmd/prism

install: ## Install the binary
	go install ./cmd/prism

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

run: ## Run the application
	go run ./cmd/prism

deps: ## Install dependencies
	go mod download
	go mod tidy

build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 go build -o bin/prism-linux-amd64 ./cmd/prism
	GOOS=linux GOARCH=arm64 go build -o bin/prism-linux-arm64 ./cmd/prism
	GOOS=darwin GOARCH=amd64 go build -o bin/prism-darwin-amd64 ./cmd/prism
	GOOS=darwin GOARCH=arm64 go build -o bin/prism-darwin-arm64 ./cmd/prism
	GOOS=windows GOARCH=amd64 go build -o bin/prism-windows-amd64.exe ./cmd/prism

.DEFAULT_GOAL := help
