# --- 1. Aşama: Build ---
# Go'nun resmi imajını temel alarak bir build ortamı oluştur
# go.mod ile uyumlu olacak şekilde Go versiyonunu 1.24'e yükseltiyoruz.
FROM golang:1.24-alpine AS builder

# Çalışma dizinini ayarla
WORKDIR /app

# git'i yükle (go mod download için gerekli olabilir)
RUN apk add --no-cache git

# go.mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./
# Bağımlılıkları indir
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

# API servisini derle. Diğer servisler (gateway vs.) için de benzer satırlar eklenebilir.
# CGO_ENABLED=0, statik bir binary oluşturmak için önemli.
# -o /app/api, derlenmiş çıktının adını ve yerini belirtir.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api ./cmd/api


# --- 2. Aşama: Final ---
# Minimal bir Alpine imajını temel alarak son imajı oluştur
FROM alpine:latest

# Çalışma dizinini ayarla
WORKDIR /app

# Yalnızca derlenmiş API binary'sini builder aşamasından kopyala
COPY --from=builder /app/api .

# Uygulamanın çalışacağı port'u dışarıya aç
EXPOSE 8080

# Konteyner başladığında çalıştırılacak komut
# Sadece derlenmiş binary'i çalıştırır.
CMD ["/app/api"] 