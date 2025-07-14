GO_PATH := $(shell go env GOPATH)

# Build variables
BINARY_NAME=anki-builder
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Platforms to build for
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: all build clean install uninstall release help

all: clean build

lint: check-lint dep
	golangci-lint run --timeout=5m -c .golangci.yml

check-lint:
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_PATH)/bin v1.64.8

dep:
	@go mod tidy
	@go mod download

# Build for current platform
build:
	@echo "Building for current platform..."
	go build ${LDFLAGS} -o build/${BINARY_NAME} ./cmd/app

# Build for all platforms
release:
	@echo "Building for all platforms..."
	@for platform in ${PLATFORMS}; do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		BINARY_SUFFIX=$${GOOS}_$${GOARCH}; \
		if [ "$$GOOS" = "windows" ]; then \
			BINARY_SUFFIX=$${BINARY_SUFFIX}.exe; \
		fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build ${LDFLAGS} -o build/${BINARY_NAME}_$${BINARY_SUFFIX} ./cmd/app; \
	done

# Install to system (Linux/macOS)
install: build
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	@sudo cp build/${BINARY_NAME} /usr/local/bin/
	@sudo chmod +x /usr/local/bin/${BINARY_NAME}
	@echo "Installation complete! Run '${BINARY_NAME} --help' to get started."

# Install to user directory (no sudo required)
install-user: build
	@echo "Installing ${BINARY_NAME} to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp build/${BINARY_NAME} ~/.local/bin/
	@chmod +x ~/.local/bin/${BINARY_NAME}
	@echo "Installation complete!"
	@echo "Add ~/.local/bin to your PATH if not already done:"
	@echo "export PATH=\$$PATH:~/.local/bin"
	@echo "Run '${BINARY_NAME} --help' to get started."

# Uninstall from system
uninstall:
	@echo "Uninstalling ${BINARY_NAME}..."
	@sudo rm -f /usr/local/bin/${BINARY_NAME}
	@echo "Uninstallation complete."

# Uninstall from user directory
uninstall-user:
	@echo "Uninstalling ${BINARY_NAME} from user directory..."
	@rm -f ~/.local/bin/${BINARY_NAME}
	@echo "Uninstallation complete."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/${BINARY_NAME}*
	@echo "Clean complete."

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	go fmt ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build for current platform"
	@echo "  release       - Build for all platforms (linux, darwin, windows)"
	@echo "  install       - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user  - Install to ~/.local/bin (no sudo required)"
	@echo "  uninstall     - Remove from /usr/local/bin"
	@echo "  uninstall-user- Remove from ~/.local/bin"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help" 