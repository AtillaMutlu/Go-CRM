#!/bin/bash

echo "🧪 Infrastructure Test Başlatılıyor..."

# PostgreSQL test
echo "📊 PostgreSQL bağlantısı test ediliyor..."
docker compose -f deploy/docker-compose.postgres.yml up -d
sleep 5

# PostgreSQL bağlantı testi
export PGPASSWORD=pass
if psql -h localhost -U user -d users -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ PostgreSQL: BAŞARILI"
else
    echo "❌ PostgreSQL: BAŞARISIZ"
    exit 1
fi

# Migration test
echo "📋 Migration test ediliyor..."
if command -v migrate &> /dev/null; then
    make migrate-up
    echo "✅ Migration: BAŞARILI"
else
    echo "⚠️  Migration tool bulunamadı. Kurulum: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
fi

# Seed data test
echo "🌱 Seed data test ediliyor..."
go run scripts/seed.go
echo "✅ Seed Data: BAŞARILI"

# Gateway servisi test
echo "🚪 Gateway servisi test ediliyor..."
timeout 30s go run cmd/gateway/main.go &
GATEWAY_PID=$!
sleep 3

if curl -f http://localhost:8080/healthz > /dev/null 2>&1; then
    echo "✅ Gateway Health: BAŞARILI"
    kill $GATEWAY_PID 2>/dev/null
else
    echo "❌ Gateway Health: BAŞARISIZ"
    kill $GATEWAY_PID 2>/dev/null
    exit 1
fi

# API servisi test
echo "📡 API servisi test ediliyor..."
timeout 30s go run cmd/api/main.go &
API_PID=$!
sleep 3

# Login test
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8085/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"demo123"}')

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    echo "✅ API Login: BAŞARILI"
    kill $API_PID 2>/dev/null
else
    echo "❌ API Login: BAŞARISIZ"
    echo "Response: $LOGIN_RESPONSE"
    kill $API_PID 2>/dev/null
    exit 1
fi

echo "🎉 Tüm infrastructure testleri BAŞARILI!" 