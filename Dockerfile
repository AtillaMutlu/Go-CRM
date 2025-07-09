# --- Build aşaması ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Servisleri build et
ARG SERVICE=gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/$SERVICE ./cmd/$SERVICE
# Seeder'ı build et
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/seeder ./scripts/seed.go

# --- Runtime aşaması ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl bash postgresql-client

WORKDIR /app

# Copy application
ARG SERVICE=gateway
COPY --from=builder /build/$SERVICE /app/service
COPY --from=builder /build/seeder /app/seeder
RUN chmod +x /app/service /app/seeder

# Copy migrations and init script
COPY migrations /app/migrations
COPY scripts/docker-init.sh /app/init.sh
RUN chmod +x /app/init.sh

# Install migrate tool by downloading binary
RUN M_VERSION="v4.17.1" && \
    apk --no-cache add wget && \
    wget "https://github.com/golang-migrate/migrate/releases/download/${M_VERSION}/migrate.linux-amd64.tar.gz" -O migrate.tar.gz && \
    tar -xvf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    rm migrate.tar.gz && \
    apk del wget

# Health check için curl gerekli
EXPOSE 8080 8085

# Use init script as entrypoint
ENTRYPOINT ["/app/init.sh", "/app/service"] 