.PHONY: help install build run test clean docker-build docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies for all services
	@echo "Installing dependencies..."
	cd order-service && go mod download
	cd payment-service && go mod download
	cd notification-service && go mod download
	@echo "Dependencies installed!"

build: ## Build all services
	@echo "Building services..."
	cd order-service && go build -o ../bin/order-service main.go
	cd payment-service && go build -o ../bin/payment-service main.go
	cd notification-service && go build -o ../bin/notification-service main.go
	@echo "Build complete! Binaries in ./bin/"

run-order: ## Run order service
	cd order-service && go run main.go

run-payment: ## Run payment service
	cd payment-service && go run main.go

run-notification: ## Run notification service
	cd notification-service && go run main.go

test: ## Run tests for all services
	@echo "Running tests..."
	cd order-service && go test ./...
	cd payment-service && go test ./...
	cd notification-service && go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf logs/
	cd order-service && go clean
	cd payment-service && go clean
	cd notification-service && go clean
	@echo "Clean complete!"

docker-build: ## Build Docker images
	docker-compose build

docker-up: ## Start all services with Docker
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

fmt: ## Format code
	cd order-service && go fmt ./...
	cd payment-service && go fmt ./...
	cd notification-service && go fmt ./...
	cd shared && go fmt ./...

vet: ## Run go vet
	cd order-service && go vet ./...
	cd payment-service && go vet ./...
	cd notification-service && go vet ./...
	cd shared && go vet ./...
