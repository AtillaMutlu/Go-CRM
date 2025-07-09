#!/bin/bash

set -e

echo "Starting Docker initialization..."

# Environment variables
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-user}
DB_PASSWORD=${DB_PASSWORD:-pass}
DB_NAME=${DB_NAME:-users}

DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Run migrations
echo "Running migrations..."
/app/migrate -path /app/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
echo "Migrations completed!"


# Seed data
echo "Seeding data..."
/app/seeder
echo "Seed data is ready!"

echo "Docker initialization complete!"
echo "Starting application..."
exec /app/service 