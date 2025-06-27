# Afrikpay Gateway - Makefile
# ================================

.PHONY: help build test clean dev stop logs docker-build docker-up docker-down

# Variables
SERVICES := auth crud temporal client
DOCKER_COMPOSE := docker-compose
GO_VERSION := 1.21

# Default target
help: ## Show this help message
	@echo "Afrikpay Gateway - Available Commands:"
	@echo "======================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development Commands
dev: ## Start all services in development mode
	@echo "🚀 Starting Afrikpay Gateway in development mode..."
	$(DOCKER_COMPOSE) up --build

dev-detached: ## Start all services in background
	@echo "🚀 Starting Afrikpay Gateway in background..."
	$(DOCKER_COMPOSE) up -d --build

stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	$(DOCKER_COMPOSE) down

restart: ## Restart all services
	@echo "🔄 Restarting all services..."
	$(DOCKER_COMPOSE) restart

logs: ## Show logs for all services
	$(DOCKER_COMPOSE) logs -f

logs-auth: ## Show logs for auth service
	$(DOCKER_COMPOSE) logs -f auth

logs-crud: ## Show logs for crud service
	$(DOCKER_COMPOSE) logs -f crud

logs-temporal: ## Show logs for temporal service
	$(DOCKER_COMPOSE) logs -f temporal

# Build Commands
build: ## Build all services
	@echo "🔨 Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd services/$$service && go build -o ../../bin/$$service ./cmd/main.go && cd ../..; \
	done

build-auth: ## Build auth service only
	@echo "🔨 Building auth service..."
	cd services/auth && go build -o ../../bin/auth ./cmd/main.go

build-crud: ## Build crud service only
	@echo "🔨 Building crud service..."
	cd services/crud && go build -o ../../bin/crud ./cmd/main.go

build-temporal: ## Build temporal service only
	@echo "🔨 Building temporal service..."
	cd services/temporal && go build -o ../../bin/temporal ./cmd/main.go

# Test Commands
test: ## Run all tests
	@echo "🧪 Running all tests..."
	go work sync
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd services/$$service && go test -v ./... && cd ../..; \
	done
	cd shared && go test -v ./...

test-auth: ## Run auth service tests
	@echo "🧪 Testing auth service..."
	cd services/auth && go test -v ./...

test-crud: ## Run crud service tests
	@echo "🧪 Testing crud service..."
	cd services/crud && go test -v ./...

test-temporal: ## Run temporal service tests
	@echo "🧪 Testing temporal service..."
	cd services/temporal && go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@for service in $(SERVICES); do \
		echo "Coverage for $$service..."; \
		cd services/$$service && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html && cd ../..; \
	done

# Docker Commands
docker-build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up: ## Start services with Docker
	@echo "🐳 Starting services with Docker..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop and remove Docker containers
	@echo "🐳 Stopping Docker containers..."
	$(DOCKER_COMPOSE) down -v

docker-clean: ## Clean Docker images and containers
	@echo "🧹 Cleaning Docker resources..."
	docker system prune -f
	docker volume prune -f

# Database Commands
db-up: ## Start databases only
	@echo "🗄️ Starting databases..."
	$(DOCKER_COMPOSE) up -d mongodb postgresql

db-down: ## Stop databases
	@echo "🗄️ Stopping databases..."
	$(DOCKER_COMPOSE) stop mongodb postgresql

# Utility Commands
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf services/*/coverage.out
	rm -rf services/*/coverage.html
	go clean -cache

deps: ## Download and tidy dependencies
	@echo "📦 Managing dependencies..."
	go work sync
	@for service in $(SERVICES); do \
		echo "Tidying $$service..."; \
		cd services/$$service && go mod tidy && cd ../..; \
	done
	cd shared && go mod tidy

fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@for service in $(SERVICES); do \
		cd services/$$service && go fmt ./... && cd ../..; \
	done
	cd shared && go fmt ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd services/$$service && golangci-lint run && cd ../..; \
	done

# Setup Commands
setup: ## Initial project setup
	@echo "⚙️ Setting up project..."
	mkdir -p bin logs
	$(MAKE) deps
	$(MAKE) docker-build

# Health Check
health: ## Check services health
	@echo "🏥 Checking services health..."
	@curl -f http://localhost:8001/health || echo "Auth service down"
	@curl -f http://localhost:8002/health || echo "CRUD service down"
	@curl -f http://localhost:8003/health || echo "Temporal service down"
