.PHONY: all build clean test run install help

BINARY_NAME=blitz
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

VERSION=2.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

all: test build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete!"

## build-all: Build for all platforms
build-all: clean
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_UNIX)_amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_UNIX)_arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_WINDOWS) .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_DARWIN)_amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_DARWIN)_arm64 .
	@echo "Cross-compilation complete! Binaries are in build/"

## run: Run the application (requires -url parameter)
run: build
	./$(BINARY_NAME) $(ARGS)

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build files
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_DARWIN)
	rm -rf build/
	rm -f coverage.out coverage.html
	rm -f blitz_results.txt
	@echo "Clean complete!"

## install: Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .
	@echo "Install complete!"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	@echo "Dependencies ready!"

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	go mod tidy

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete!"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed!"

## help: Show this help message
help:
	@echo "Blitz - Makefile commands:"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /'

.DEFAULT_GOAL := help
