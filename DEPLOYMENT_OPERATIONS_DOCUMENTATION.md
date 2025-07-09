# Dağıtım ve Operasyon Dokümantasyonu

## İçindekiler
1. Genel Bakış
2. Ön Gereksinimler
3. Yerel Geliştirme Kurulumu
4. Docker ile Dağıtım
5. Production Dağıtımı
6. Konfigürasyon Yönetimi
7. Veritabanı Yönetimi
8. İzleme ve Loglama
9. Güvenlik
10. Yedekleme ve Kurtarma
11. Sorun Giderme
12. Bakım
13. Blue/Green Deployment
14. Load Balancer (Nginx)
15. Primary/Replica PostgreSQL
16. Otomatik Migration
17. Merkezi Loglama
18. Rate Limiting ve Cache (Redis)
19. Event Publish (Kafka)
20. Prometheus & Jaeger İzleme
21. JWT ile Güvenlik
22. Swagger/OpenAPI
23. CI/CD Pipeline

---

## 1. Genel Bakış
Bu doküman, Go tabanlı mikroservis CRM uygulamasının farklı ortamlarda dağıtımı, işletimi ve bakımı için kapsamlı rehber sunar. Tüm mimari, kurumsal seviyede ölçeklenebilirlik, güvenlik ve sürdürülebilirlik hedefiyle tasarlanmıştır.

### Temel Mimarî Bileşenler
- Frontend: Nginx ile sunulan statik dosyalar
- API Gateway: Kimlik doğrulama, rate limiting, merkezi loglama, izleme
- Mikroservisler: Her iş alanı için bağımsız Go servisi
- Veritabanı: PostgreSQL (primary/replica mimarisi)
- Mesajlaşma: Kafka
- Cache & Rate Limiting: Redis
- İzleme: Prometheus, Jaeger
- Loglama: Logrus (JSON, request-id/trace-id)
- API Dokümantasyonu: Swagger/OpenAPI
- CI/CD: GitHub Actions

---

## 2. Ön Gereksinimler
- Docker 20.10+
- Docker Compose 2.0+
- Make
- Git
- (Opsiyonel) Docker Desktop, kubectl, helm

Donanım:
- Minimum: 2 CPU, 4GB RAM, 20GB disk
- Önerilen: 4+ CPU, 8GB+ RAM, 50GB+ SSD

Ağ:
- 3000: Frontend
- 8080: API Gateway
- 8085: API Servisi
- 5432: PostgreSQL
- 6379: Redis
- 9092: Kafka
- 9090: Prometheus
- 16686: Jaeger

---

## 3. Yerel Geliştirme Kurulumu
1. Repoyu klonlayın:
```
git clone <repo-url>
cd Go-CRM
```
2. Servisleri başlatın:
```
make docker-up
```
3. Uygulamaya erişin: http://localhost:3000
4. Servisleri durdurun:
```
make docker-down
```

---

## 4. Docker ile Dağıtım
Tüm servisler docker-compose ile ayağa kaldırılır. Blue/Green deployment, Nginx load balancer, PostgreSQL primary/replica, Redis, Kafka, Prometheus, Jaeger, Grafana gibi bileşenler desteklenir.

---

## 5. Production Dağıtımı
- Blue/Green deployment ile sıfır kesintiyle güncelleme
- Nginx ile yük dengeleme ve SSL sonlandırma
- PostgreSQL primary/replica ile yüksek erişilebilirlik
- Redis ile rate limiting ve cache
- Kafka ile event publish/subscribe
- Prometheus ve Jaeger ile izleme ve tracing
- CI/CD pipeline ile otomatik test, build ve deployment

---

## 6. Konfigürasyon Yönetimi
- Ortam değişkenleri .env dosyaları ile yönetilir
- Hassas bilgiler için Vault veya Docker secrets entegrasyonu önerilir

---

## 7. Veritabanı Yönetimi
- migrations/ klasöründe SQL migration dosyaları tutulur
- make docker-up ile migration’lar otomatik uygulanır
- Production’da primary/replica mimarisi ile okuma/yazma ayrımı yapılabilir

---

## 8. İzleme ve Loglama
- Prometheus ile metrik toplama
- Jaeger ile distributed tracing
- Logrus ile JSON formatında merkezi loglama (request-id/trace-id ile)
- ELK veya Loki entegrasyonu için hazır yapı

---

## 9. Güvenlik
- JWT ile authentication
- Rate limiting (Redis tabanlı)
- CORS ve güvenli header’lar
- Hassas bilgiler için secret management

---

## 10. Yedekleme ve Kurtarma
- PostgreSQL volume’ları düzenli yedeklenmeli
- Redis ve Kafka için de yedekleme stratejisi uygulanmalı

---

## 11. Sorun Giderme
- make logs ile tüm servis loglarını görebilirsiniz
- docker-compose logs <servis> ile spesifik servis logları alınabilir
- make db-connect ile veritabanına bağlanabilirsiniz

---

## 12. Bakım
- Tüm modüller için test ve dokümantasyon zorunludur
- Kod ve dokümantasyonda ikon, emoji veya süsleme kullanılmaz
- Her commit ve PR Türkçe açıklamalı olmalıdır

---

## 13. Blue/Green Deployment
- Blue/Green deployment ile yeni versiyonlar canlıya alınırken eski versiyon anında geri döndürülebilir
- docker-compose.db-bluegreen.yml örneği ile iki ayrı ortam paralel yönetilebilir

---

## 14. Load Balancer (Nginx)
- Nginx reverse proxy olarak çalışır
- Tüm HTTP trafiğini karşılar, frontend ve API isteklerini ilgili servislere yönlendirir
- SSL sonlandırma ve healthcheck desteği vardır

---

## 15. Primary/Replica PostgreSQL
- Okuma/yazma yükleri ayrılabilir
- docker-compose.db-bluegreen.yml ile primary/replica kurulumu örneklenmiştir
- Yüksek erişilebilirlik ve performans sağlar

---

## 16. Otomatik Migration
- migrations/ klasörüne eklenen SQL dosyaları make docker-up ile otomatik uygulanır
- Production’da migration işlemleri CI/CD pipeline’ında otomatikleştirilebilir

---

## 17. Merkezi Loglama
- Tüm loglar JSON formatında, request-id/trace-id ile tutulur
- Logrus kullanılır, merkezi log sistemlerine (ELK, Loki) kolayca entegre edilebilir

---

## 18. Rate Limiting ve Cache (Redis)
- Login ve create endpoint’lerinde Redis tabanlı rate limiting (5 istek/dk/IP)
- Müşteri listeleme gibi işlemlerde Redis cache kullanılır

---

## 19. Event Publish (Kafka)
- Müşteri oluşturulunca Kafka ile event publish edilir
- Servisler arası asenkron iletişim için altyapı hazırdır

---

## 20. Prometheus & Jaeger İzleme
- Prometheus ile metrik toplama ve izleme
- Jaeger ile distributed tracing
- Grafana ile görselleştirme mümkündür

---

## 21. JWT ile Güvenlik
- Tüm korumalı endpoint’lerde JWT doğrulama zorunludur
- Şifreler bcrypt ile hashlenir

---

## 22. Swagger/OpenAPI
- Handler fonksiyonlarında açıklamalar ile otomatik API dokümantasyonu
- /swagger endpoint’i ile erişim
- Swaggo/swag ile Go kodundan otomatik üretim (Windows’ta PATH ayarı gerekebilir)

---

## 23. CI/CD Pipeline
- GitHub Actions ile otomatik test, build ve deployment
- Otomatik migration ve test coverage raporu

---

Tüm detaylar ve örnekler için ilgili YAML ve kod dosyalarına bakınız. Her türlü katkı ve öneri için proje kurallarına uyunuz.