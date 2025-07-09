# Ortak değişkenler
SERVICE ?= gateway
BINARY = ./build/$(SERVICE)

MIGRATE ?= $(HOME)/go/bin/migrate
MIGRATIONS_PATH = migrations
DB_URL = postgres://user:pass@localhost:5432/users?sslmode=disable

.PHONY: all build run lint test docker migrate-up migrate-down

all: build

build:
	go build -o $(BINARY) ./cmd/$(SERVICE)

run:
	go run ./cmd/$(SERVICE)

lint:
	golangci-lint run ./...

test:
	go test ./...

docker:
	docker build -t $(SERVICE):dev --build-arg SERVICE=$(SERVICE) --build-arg ENTRY=$(SERVICE) . 

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1 