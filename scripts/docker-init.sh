#!/bin/bash

set -e

echo "🚀 Docker initialization başlatılıyor..."

# Environment variables
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-user}
DB_PASSWORD=${DB_PASSWORD:-pass}
DB_NAME=${DB_NAME:-users}

DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

echo "📊 Veritabanı bağlantısı bekleniyor: $DB_HOST:$DB_PORT"

# Wait for database
max_attempts=30
attempt=1

while [ $attempt -le $max_attempts ]; do
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
        echo "✅ Veritabanı hazır!"
        break
    fi
    
    echo "⏳ Veritabanı bekleniyor... ($attempt/$max_attempts)"
    sleep 2
    attempt=$((attempt + 1))
done

if [ $attempt -gt $max_attempts ]; then
    echo "❌ Veritabanı bağlantısı timeout!"
    exit 1
fi

# Run migrations if migrate tool exists
if command -v migrate &> /dev/null; then
    echo "📋 Migration'lar çalıştırılıyor..."
    # Dirty state'i temizle ve migration'ı zorla
    migrate -path /app/migrations -database "$DATABASE_URL" force 1
    migrate -path /app/migrations -database "$DATABASE_URL" up
    echo "✅ Migration'lar tamamlandı!"
else
    echo "⚠️  Migrate tool bulunamadı, migration atlanıyor."
fi

# Run seeder if seed script exists
if [ -f "/app/seeder" ]; then
    echo "🌱 Seed data oluşturuluyor..."
    /app/seeder
    echo "✅ Seed data hazır!"
fi

echo "🎉 Docker initialization tamamlandı!"

# Start the main application
echo "🚀 Uygulama başlatılıyor..."
exec "$@" 