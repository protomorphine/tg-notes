.PHONY: help build run clean docker test lint deps generate

# ====================================================================================
# VARIABLES
# ====================================================================================

BINARY_NAME=tg-notes
CONFIG_PATH=config/local.yaml

# ====================================================================================
# HELP
# ====================================================================================

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  help:          Show this help message"
	@echo "  deps:          Install dependencies"
	@echo "  generate:      Generate mocks"
	@echo "  lint:          Run linter"
	@echo "  test:          Run tests"
	@echo "  test-verbose:  Run tests in verbose mode"
	@echo "  build:         Build the project"
	@echo "  run:           Run the project"
	@echo "  clean:         Clean up build artifacts"
	@echo "  docker:        Build a Docker image"


# ====================================================================================
# COMMANDS
# ====================================================================================

deps:
	@echo "Installing dependencies..."
	@go install github.com/vektra/mockery/v2@latest

generate:
	@echo "Generating mocks..."
	@mockery

lint:
ifeq (, $(shell which golangci-lint))
	@echo "golangci-lint is not installed. Please install it by running:"
	@echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
else
	@echo "Running linter..."
	@golangci-lint run
endif

test:
	@echo "Running tests..."
	@go test ./...

test-verbose:
	@echo "Running tests in verbose mode..."
	@go test -v ./...

build:
	@echo "Building the project..."
	@go build -o ${BINARY_NAME} main.go bot.go logger.go

run:
	@echo "Running the project..."
	@go run main.go bot.go logger.go --config ${CONFIG_PATH}

clean:
	@echo "Cleaning up build artifacts..."
	@rm -f ${BINARY_NAME}

docker:
	@echo "Building a Docker image..."
	@docker build -t ${BINARY_NAME} .
