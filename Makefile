.PHONY: build run test clean install lint docker-build

BINARY_NAME=openshift-mcp
BUILD_DIR=./build
CMD_DIR=./cmd/server

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go

run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(CMD_DIR)/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linter..."
	@golangci-lint run

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest -f build/Dockerfile .
