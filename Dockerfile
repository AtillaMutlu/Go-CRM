# --- Build aşaması ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Servis adı argümanı
ARG SERVICE=gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/$SERVICE ./cmd/$SERVICE

# --- Runtime aşaması ---
FROM gcr.io/distroless/static:nonroot
ARG SERVICE=gateway
COPY --from=builder /build/$SERVICE /$SERVICE
USER nonroot:nonroot
ENTRYPOINT ["/gateway"] 