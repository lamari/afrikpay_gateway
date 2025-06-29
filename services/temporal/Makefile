# Makefile for Afrikpay Gateway

.PHONY: help test coverage docker-up docker-down lint fmt e2e-newman

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  test        Run all Go unit tests"
	@echo "  coverage    Run tests with coverage and generate HTML report"
	@echo "  docker-up   Start all services with docker-compose"
	@echo "  docker-down Stop all services with docker-compose"
	@echo "  lint        Run golangci-lint (if installed)"
	@echo "  fmt         Run go fmt on all Go files"
	@echo "  e2e-newman  Run E2E API tests with Newman (Postman CLI)"

# Run unit tests

test:
	go test ./... -v

# Run tests with coverage and open HTML report
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker Compose commands
docker-build:
	docker-compose down && docker-compose up --build

docker-up:
	docker-compose down && docker-compose up

docker-down:
	docker-compose down

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

# Format Go code
fmt:
	go fmt ./...

# Build API app
build-api:
	go build -o ./bin/api ./cmd/api

# Build Worker app
build-worker:
	go build -o ./bin/worker ./cmd/worker

build:
	make build-api && make build-worker

# Run E2E API tests with Newman (Postman CLI)
e2e-newman:
	newman run e2e/temporal_schema_e2e.postman_collection.json --env-var base_url=http://localhost:8080
