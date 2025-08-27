# Log Ingestion Server Makefile

.PHONY: build run test clean docker-build docker-run migrate-up migrate-down deps lint format

# Variables
APP_NAME=log-ingestion-server
VERSION=1.0.0
BUILD_DIR=build
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# Default target
all: deps lint format test build

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.VERSION=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) .

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Lint the code
lint:
	golangci-lint run

# Format the code
format:
	go fmt ./...
	goimports -w .

# Run database migrations up
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

# Create a new migration
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run:
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	cp config.example.env .env
	@echo "Please edit .env file with your configuration"
	@echo "Run 'make migrate-up' to setup database"

# Production build
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.VERSION=$(VERSION) -w -s" -o $(BUILD_DIR)/$(APP_NAME) .

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	@echo "API documentation would be generated here"

# Health check
health:
	curl -f http://localhost:8080/health || exit 1

# Load test (requires wrk or similar tool)
load-test:
	@echo "Running load test..."
	@echo "Install 'wrk' tool first: brew install wrk"
	wrk -t12 -c400 -d30s --header "X-API-Key: your-api-key" http://localhost:8080/health

# Show help
help:
	@echo "Available targets:"
	@echo "  deps          - Install dependencies"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Lint the code"
	@echo "  format        - Format the code"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  migrate-create - Create a new migration (use: make migrate-create name=migration_name)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-up     - Start with docker-compose"
	@echo "  docker-down   - Stop docker-compose"
	@echo "  dev-setup     - Setup development environment"
	@echo "  build-prod    - Build for production"
	@echo "  install-tools - Install development tools"
	@echo "  health        - Check server health"
	@echo "  help          - Show this help message"
