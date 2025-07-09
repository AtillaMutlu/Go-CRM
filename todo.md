# Mikroservis Örnekleri ve Geliştirilecekler

## Yapılanlar

- **Gateway Servisi**
  - JWT doğrulama (Keycloak JWKS ile)
  - Rate limit (IP başına 10 istek/dk)
  - IP allowlist (127.0.0.1, ::1)

- **Auth Servisi**
  - Healthz ve info endpointleri

- **User Servisi (gRPC)**
  - Protobuf tanımı ve Go kodu
  - Temel gRPC server iskeleti

- **PostgreSQL**
  - Primary ve read-replica docker-compose

- **Kafka**
  - Kafka + Zookeeper docker-compose

- **Notification Servisi**
  - Kafka'dan notification.command topic'ini dinliyor
  - Gelen mesajı logluyor

- **Audit Servisi**
  - Kafka'dan audit.raw topic'ini dinliyor
  - Mesajı immudb'ye append ediyor

## Geliştirilecekler (örnekler)

- **Dead-letter queue (notification.dlq)**
  - Başarısız bildirimler için ayrı Kafka kuyruğu

- **Retry policy (exponential backoff, max 5)**
  - Bildirim gönderimi başarısız olursa otomatik tekrar deneme

- **Testler**
  - Unit ve integration test örnekleri

- **SMTP, SMS, WebSocket adapterleri**
  - Notification servisinde farklı kanal entegrasyonları

- **Audit için Merkle root S3'e dump**
  - immudb Merkle root'unu günlük olarak S3'e kaydetme

- **Observability**
  - Prometheus, Grafana, Jaeger entegrasyonu

- **Helm chart, Argo Rollouts**
  - K8s deployment ve canary örnekleri 