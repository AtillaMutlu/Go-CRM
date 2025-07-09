# Go CRM Mikroservisleri - API Dokümantasyonu

## İçindekiler
1. Genel Bakış
2. Mimari
3. Kimlik Doğrulama (JWT)
4. REST API Endpoint’leri
5. gRPC Servisleri
6. Rate Limiting ve Cache
7. Event Publish (Kafka)
8. Merkezi Loglama
9. İzleme (Prometheus & Jaeger)
10. Swagger/OpenAPI
11. CI/CD Pipeline

---

## 1. Genel Bakış
Bu doküman, Go tabanlı mikroservis CRM uygulamasının API katmanını ve mimari detaylarını açıklar. Tüm sistem, kurumsal seviyede ölçeklenebilirlik, güvenlik ve sürdürülebilirlik hedefiyle tasarlanmıştır.

### Temel Teknolojiler
- Backend: Go
- API: RESTful + gRPC (hibrit iletişim)
- Kimlik Doğrulama: JWT
- Rate Limiting & Cache: Redis
- Mesajlaşma: Kafka
- İzleme: Prometheus, Jaeger
- Loglama: Logrus (JSON, request-id/trace-id)
- API Dokümantasyonu: Swagger/OpenAPI
- CI/CD: GitHub Actions

---

## 2. Mimari
Her servis kendi konteynerinde çalışır. Nginx yük dengeleyici olarak görev yapar. API Gateway, kimlik doğrulama, rate limiting, merkezi loglama ve izleme işlevlerini üstlenir. Mikroservisler arası iletişim gRPC ve Kafka ile sağlanır. PostgreSQL primary/replica mimarisi ile yüksek erişilebilirlik sağlanır.

---

## 3. Kimlik Doğrulama (JWT)
Tüm korumalı endpoint’lerde JWT doğrulama zorunludur. Giriş işlemi sonrası JWT token döner. Token 24 saat geçerlidir. Authorization header’ı ile gönderilmelidir.

---

## 4. REST API Endpoint’leri
### Müşteri Yönetimi
- GET /api/customers : Tüm müşterileri listeler (JWT zorunlu)
- POST /api/customers : Yeni müşteri oluşturur (JWT zorunlu, rate limitli)
- PUT /api/customers/{id} : Müşteri günceller (JWT zorunlu)
- DELETE /api/customers/{id} : Müşteri siler (JWT zorunlu)

### Kullanıcı Yönetimi
- POST /api/login : Giriş ve JWT token alma (rate limitli)
- POST /api/register : Yeni kullanıcı oluşturma

### Diğer
- GET /api/healthz : Healthcheck endpoint’i
- GET /api/metrics : Prometheus metrikleri
- GET /swagger : Swagger/OpenAPI dokümantasyonu

---

## 5. gRPC Servisleri
Mikroservisler arası yüksek performanslı iletişim için gRPC protokolleri kullanılır. Protobuf tanımları proto/ klasöründedir.

---

## 6. Rate Limiting ve Cache
- Login ve create endpoint’lerinde Redis tabanlı rate limiting (5 istek/dk/IP)
- Müşteri listeleme gibi işlemlerde Redis cache kullanılır

---

## 7. Event Publish (Kafka)
- Müşteri oluşturulunca Kafka ile event publish edilir
- Servisler arası asenkron iletişim için altyapı hazırdır

---

## 8. Merkezi Loglama
- Tüm loglar JSON formatında, request-id/trace-id ile tutulur
- Logrus kullanılır, merkezi log sistemlerine (ELK, Loki) kolayca entegre edilebilir

---

## 9. İzleme (Prometheus & Jaeger)
- Prometheus ile metrik toplama ve izleme
- Jaeger ile distributed tracing

---

## 10. Swagger/OpenAPI
- Handler fonksiyonlarında açıklamalar ile otomatik API dokümantasyonu
- /swagger endpoint’i ile erişim

---

## 11. CI/CD Pipeline
- GitHub Actions ile otomatik test, build ve deployment
- Otomatik migration ve test coverage raporu

---

Tüm endpoint örnekleri, hata mesajları ve detaylar Swagger/OpenAPI dokümantasyonunda ve kodda bulunabilir. Her türlü katkı ve öneri için proje kurallarına uyunuz.