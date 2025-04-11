.PHONY: build clean test lint run help install release snapshot docker docs completions all

# Binary name
BINARY_NAME=swagger-to-http
# Main package path
MAIN_PACKAGE=./cmd/swagger-to-http
# Output directory
OUTPUT_DIR=bin

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Platforms to build for
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Linker flags
LDFLAGS=-ldflags "-s -w -X github.com/edgardnogueira/swagger-to-http/internal/version.Version=$(VERSION) -X github.com/edgardnogueira/swagger-to-http/internal/version.Commit=$(COMMIT) -X github.com/edgardnogueira/swagger-to-http/internal/version.BuildDate=$(BUILD_DATE)"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Build for all platforms
build-all:
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(OUTPUT_DIR)
	$(foreach platform,$(PLATFORMS),\
		$(eval GOOS=$(word 1,$(subst /, ,$(platform)))) \
		$(eval GOARCH=$(word 2,$(subst /, ,$(platform)))) \
		$(eval OUTPUT=$(if $(findstring windows,$(GOOS)),$(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe,$(OUTPUT_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH))) \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(OUTPUT) $(MAIN_PACKAGE); \
	)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)/
	@rm -rf dist/
	@rm -rf completions/
	@rm -rf manpages/
	@go clean -cache -testcache -modcache

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

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...

# Install linting tools
lint-tools:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

# Run linting
lint:
	@echo "Running linters..."
	@golangci-lint run ./...

# Run the application
run:
	@go run $(LDFLAGS) $(MAIN_PACKAGE) $(ARGS)

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

# Build Docker image
docker:
	@echo "Building Docker image..."
	@docker build -t edgardnogueira/$(BINARY_NAME):$(VERSION) .

# Run Docker image
docker-run:
	@echo "Running Docker image..."
	@docker run --rm -it edgardnogueira/$(BINARY_NAME):$(VERSION) $(ARGS)

# Build documentation
docs:
	@echo "Building documentation..."
	@mkdir -p docs
	@go run $(MAIN_PACKAGE) generate-docs docs

# Create shell completions
completions:
	@echo "Generating shell completions..."
	@mkdir -p completions
	@mkdir -p manpages
	@go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@$(OUTPUT_DIR)/$(BINARY_NAME) completion bash > completions/$(BINARY_NAME).bash
	@$(OUTPUT_DIR)/$(BINARY_NAME) completion zsh > completions/$(BINARY_NAME).zsh
	@$(OUTPUT_DIR)/$(BINARY_NAME) completion fish > completions/$(BINARY_NAME).fish
	@$(OUTPUT_DIR)/$(BINARY_NAME) man > manpages/$(BINARY_NAME).1
	@gzip -f manpages/$(BINARY_NAME).1

# Tidy up dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	@go mod verify

# Generate licenses report
licenses:
	@echo "Generating licenses report..."
	@go list -m -json all | go-licenses report -t NOTICE > THIRD_PARTY_LICENSES.txt

# All build steps for CI
all: tidy verify lint test build

# Help command
help:
	@echo "Available commands:"
	@echo "  all         - Run tidy, verify, lint, test, and build"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for all platforms"
	@echo "  clean       - Clean build artifacts"
	@echo "  completions - Generate shell completions"
	@echo "  cover       - Run tests with coverage"
	@echo "  docker      - Build Docker image"
	@echo "  docker-run  - Run Docker image"
	@echo "  docs        - Build documentation"
	@echo "  install     - Install the application"
	@echo "  licenses    - Generate licenses report"
	@echo "  lint        - Run linters"
	@echo "  lint-tools  - Install linting tools"
	@echo "  release     - Create a new release with GoReleaser"
	@echo "  run         - Run the application"
	@echo "  snapshot    - Create a snapshot release for testing"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race detection"
	@echo "  tidy        - Tidy up dependencies"
	@echo "  verify      - Verify dependencies"
	@echo "  help        - Show this help message"

# Default target
default: build
