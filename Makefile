# Ortak değişkenler
SERVICE ?= gateway
BINARY = ./build/$(SERVICE)

MIGRATE ?= $(HOME)/go/bin/migrate
MIGRATIONS_PATH = migrations
DB_URL = postgres://user:pass@localhost:5432/users?sslmode=disable

.PHONY: all build run lint test docker migrate-up migrate-down test-unit test-integration test-e2e test-all test-infrastructure docker-up docker-up-d docker-down docker-logs docker-status docker-rebuild docker-clean frontend

all: build

build:
	go build -o $(BINARY) ./cmd/$(SERVICE)

run:
	go run ./cmd/$(SERVICE)

lint:
	golangci-lint run ./...

test:
	go test ./...

# Test komutları
test-infrastructure:
	@echo "🧪 Infrastructure testleri çalıştırılıyor..."
	@bash scripts/test-infrastructure.sh

test-unit:
	@echo "🧪 Unit testler çalıştırılıyor..."
	@go test -v ./tests/unit/...

test-integration:
	@echo "🧪 Integration testler çalıştırılıyor..."
	@export TEST_DB_URL="postgres://user:pass@localhost:5432/users?sslmode=disable" && go test -v ./tests/integration/...

test-e2e:
	@echo "🧪 E2E testler çalıştırılıyor..."
	@go test -v ./tests/e2e/...

test-all:
	@echo "🧪 Tüm testler çalıştırılıyor..."
	@bash scripts/run-all-tests.sh

test-performance:
	@echo "🚀 Performance testleri çalıştırılıyor..."
	@go test -bench=. -benchmem ./tests/unit/...

# Test coverage
test-coverage:
	@echo "📊 Test coverage raporu oluşturuluyor..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage raporu coverage.html dosyasında hazır!"

docker:
	docker build -t $(SERVICE):dev --build-arg SERVICE=$(SERVICE) --build-arg ENTRY=$(SERVICE) . 

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

# Test migration'ları
migrate-test-up:
	$(MIGRATE) -path migrations/test -database "$(DB_URL)" up

migrate-test-down:
	$(MIGRATE) -path migrations/test -database "$(DB_URL)" down

# Docker komutları
docker-up:
	@echo "🚀 Docker Compose ile tüm servisleri başlatıyor..."
	docker-compose up --build

docker-up-d:
	@echo "🚀 Docker Compose (detached mode)..."
	docker-compose up --build -d

docker-down:
	@echo "🛑 Docker servisleri durduruluyor..."
	docker-compose down

docker-logs:
	@echo "📋 Docker logları..."
	docker-compose logs -f

docker-status:
	@echo "📊 Docker servis durumları..."
	docker-compose ps

docker-rebuild:
	@echo "🔄 Docker images yeniden oluşturuluyor..."
	docker-compose build --no-cache

docker-clean:
	@echo "🧹 Docker temizliği..."
	docker-compose down -v
	docker system prune -f

# Frontend'i çalıştır
frontend:
	@echo "🌐 Frontend http://localhost:3000 adresinde çalışıyor..."
	@echo "📡 API Proxy: http://localhost:3000/api"
	@echo "🔑 Login: demo@example.com / demo123"

# Cleanup
clean:
	rm -f ./build/*
	rm -f coverage.out coverage.html 