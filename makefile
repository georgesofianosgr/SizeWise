.PHONY: build test clean run deps help

# Binary name for the server
BINARY_NAME=sizewise
TAG=latest

# Default Go command
GO=go

# Ensure the build directory exists and build the project
build:
	mkdir -p build
	$(GO) build -o build/$(BINARY_NAME) ./cmd/server

docker-build:
	docker build -t sizewise:$(TAG) .

# Run tests for all modules
test:
	$(GO) test -v ./pkg/...

# Clean up binaries and the build directory
clean:
	$(GO) clean
	rm -rf build

# Run the server
run: 
	APP_ENV="development" ./build/$(BINARY_NAME) -serve

# Install dependencies
deps:
	$(GO) mod tidy
	$(GO) mod download

# Help command to list available commands
help:
	@echo "Available commands:"
	@echo "  build       Build the server binary in the build directory"
	@echo "  test        Run tests"
	@echo "  clean       Clean up binaries and the build directory"
	@echo "  run         Build and run the server"
	@echo "  deps        Install dependencies"
	@echo "  help        Display this help message"
