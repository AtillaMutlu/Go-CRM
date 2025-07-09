# --- Build aşaması ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Servis adı argümanı
ARG SERVICE=gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/$SERVICE ./cmd/$SERVICE

# --- Runtime aşaması ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl bash postgresql-client go

WORKDIR /app

# Copy application
ARG SERVICE=gateway
COPY --from=builder /build/$SERVICE /app/service
RUN chmod +x /app/service

# Copy source files for seeder
COPY . /app/

# Copy init script
COPY scripts/docker-init.sh /app/init.sh
RUN chmod +x /app/init.sh

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
ENV PATH="/root/go/bin:${PATH}"

# Health check için curl gerekli
EXPOSE 8080 8085

# Use init script as entrypoint
ENTRYPOINT ["/app/init.sh", "/app/service"] 