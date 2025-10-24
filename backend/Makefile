.PHONY: run build test clean docker-up docker-down migrate help

# Variables
APP_NAME=daybook-backend
DOCKER_COMPOSE=docker-compose

# Help command
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application locally"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-up    - Start all services with Docker Compose"
	@echo "  make docker-down  - Stop all services"
	@echo "  make docker-logs  - View Docker logs"
	@echo "  make docker-build - Build Docker image"

# Run the application locally
run:
	go run main.go

# Build the application
build:
	go build -o bin/$(APP_NAME) main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Start all services with Docker Compose
docker-up:
	$(DOCKER_COMPOSE) up -d

# Stop all services
docker-down:
	$(DOCKER_COMPOSE) down

# View Docker logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Build Docker image
docker-build:
	$(DOCKER_COMPOSE) build

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...
