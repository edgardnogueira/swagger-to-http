.PHONY: build clean test lint run help install

# Binary name
BINARY_NAME=swagger-to-http
# Main package path
MAIN_PACKAGE=./cmd/swagger-to-http

# Version information
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags
LDFLAGS=-ldflags "-s -w -X github.com/edgardnogueira/swagger-to-http/internal/version.Version=$(VERSION) -X github.com/edgardnogueira/swagger-to-http/internal/version.Commit=$(COMMIT) -X github.com/edgardnogueira/swagger-to-http/internal/version.BuildDate=$(BUILD_DATE)"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -rf dist/
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
cover:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Install linting tools
lint-tools:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.54.2

# Run linting
lint:
	@echo "Running linters..."
	@golangci-lint run ./...

# Run the application
run:
	@go run $(LDFLAGS) $(MAIN_PACKAGE)

# Install the application
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(MAIN_PACKAGE)

# Create a new release with GoReleaser
release:
	@echo "Creating new release with GoReleaser..."
	@goreleaser release --clean

# Prepare a snapshot release for testing
snapshot:
	@echo "Creating snapshot release with GoReleaser..."
	@goreleaser release --clean --snapshot

# Tidy up dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Help command
help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  cover     - Run tests with coverage"
	@echo "  lint      - Run linters"
	@echo "  lint-tools - Install linting tools"
	@echo "  run       - Run the application"
	@echo "  install   - Install the application"
	@echo "  release   - Create a new release with GoReleaser"
	@echo "  snapshot  - Create a snapshot release for testing"
	@echo "  tidy      - Tidy up dependencies"
	@echo "  help      - Show this help message"

# Default target
default: build
