# Go Fonksiyonları ve Bileşenleri Dokümantasyonu

## İçindekiler
1. Servis Giriş Noktaları
2. Ortak Yardımcı Fonksiyonlar
3. Middleware’ler
4. Handler Fonksiyonları
5. Repository ve Service Katmanı
6. Rate Limiting ve Cache
7. Event Publish (Kafka)
8. Merkezi Loglama
9. İzleme (Prometheus & Jaeger)
10. JWT ile Güvenlik
11. Swagger/OpenAPI
12. CI/CD Pipeline

---

## 1. Servis Giriş Noktaları
- cmd/api/main.go : CRM API servisi
- cmd/gateway/main.go : API Gateway servisi
- cmd/auth-svc/main.go : Kimlik doğrulama servisi
- cmd/user-svc/main.go : Kullanıcı yönetimi servisi
- cmd/notification-svc/main.go : Bildirim servisi
- cmd/audit-svc/main.go : Audit servisi

Her servis bağımsız olarak başlatılabilir ve ölçeklenebilir.

---

## 2. Ortak Yardımcı Fonksiyonlar
- Ortam değişkeni okuma, hata yönetimi, validasyon, şifre hashleme (bcrypt), request-id üretimi gibi fonksiyonlar pkg/ altında bulunur.

---

## 3. Middleware’ler
- JWT doğrulama
- Rate limiting (Redis tabanlı)
- CORS
- Merkezi loglama (request-id/trace-id)
- Prometheus metrikleri
- Jaeger tracing

---

## 4. Handler Fonksiyonları
- Tüm endpoint’ler için ayrı handler fonksiyonları bulunur
- Swagger/OpenAPI açıklamaları ile dokümante edilir
- Hata yönetimi ve validasyon standarttır

---

## 5. Repository ve Service Katmanı
- Repository: Veritabanı işlemleri (CRUD, raw SQL, transaction)
- Service: İş mantığı, validasyon, event publish
- Her modül için ayrı repository ve service dosyası

---

## 6. Rate Limiting ve Cache
- Redis ile login ve create endpoint’lerinde rate limiting (5 istek/dk/IP)
- Müşteri listeleme gibi işlemlerde Redis cache

---

## 7. Event Publish (Kafka)
- Müşteri oluşturulunca Kafka ile event publish edilir
- Event consumer altyapısı hazırdır

---

## 8. Merkezi Loglama
- Logrus ile JSON formatında loglama
- request-id/trace-id ile merkezi log sistemlerine uygun yapı

---

## 9. İzleme (Prometheus & Jaeger)
- Prometheus ile metrik toplama
- Jaeger ile distributed tracing

---

## 10. JWT ile Güvenlik
- JWT ile authentication middleware
- Şifreler bcrypt ile hashlenir

---

## 11. Swagger/OpenAPI
- Handler fonksiyonlarında açıklamalar ile otomatik API dokümantasyonu
- /swagger endpoint’i ile erişim

---

## 12. CI/CD Pipeline
- GitHub Actions ile otomatik test, build ve deployment
- Otomatik migration ve test coverage raporu

---

Tüm fonksiyon ve bileşen örnekleri için kodda ilgili dosyalara bakınız. Her türlü katkı ve öneri için proje kurallarına uyunuz.
