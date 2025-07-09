# Makefile

.PHONY: help build-api build-gateway run-api run-gateway proto-gen test-all test-unit test-integration test-e2e docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  build-api          - Build the API service"
	@echo "  build-gateway      - Build the Gateway service"
	@echo "  run-api            - Run the API service"
	@echo "  run-gateway        - Run the Gateway service"
	@echo "  proto-gen          - Generate gRPC code from proto file"
	@echo "  test-all           - Run all tests (unit, integration, e2e)"
	@echo "  test-unit          - Run unit tests"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-e2e           - Run end-to-end tests"
	@echo "  docker-up          - Start all services with Docker Compose"
	@echo "  docker-down        - Stop all services with Docker Compose"

# Build commands
build-api:
	@echo "Building API service..."
	go build -o ./bin/api ./cmd/api

build-gateway:
	@echo "Building Gateway service..."
	go build -o ./bin/gateway ./cmd/gateway

# Run commands
run-api:
	@echo "Running API service..."
	go run ./cmd/api/main.go

run-gateway:
	@echo "Running Gateway service..."
	go run ./cmd/gateway/main.go

# Proto generation
proto-gen:
	@echo "Generating gRPC code..."
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative proto/user.proto

# Test commands
test-all:
	@echo "Running all tests..."
	./scripts/run-all-tests.sh

test-unit:
	@echo "Running unit tests..."
	go test ./tests/unit/...

test-integration:
	@echo "Running integration tests..."
	go test ./tests/integration/...

test-e2e:
	@echo "Running end-to-end tests..."
	go test ./tests/e2e/...

test-coverage:
	@echo "Test coverage raporu üretiliyor..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "coverage.html dosyasını tarayıcıda açarak detaylı raporu görebilirsiniz."

coverage-badge:
	@echo "Coverage badge (SVG) üretiliyor..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//' | xargs -I{} curl -s "https://img.shields.io/badge/coverage-{}%25-brightgreen.svg" -o coverage.svg
	@echo "coverage.svg dosyası oluşturuldu. README'ye ekleyebilirsiniz."

# Docker commands
docker-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up --build

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down 