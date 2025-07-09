# Proje Durumu ve Geliştirilecekler

## Yapılanlar

- **Mimari ve Altyapı**
  - Docker Compose ile tüm projenin (API, Gateway, Frontend, DB) tek komutla (`make docker-up`) çalıştırılması.
  - Otomatik veritabanı migration ve seeder mekanizması.
  - Nginx ile frontend sunumu ve API Gateway için reverse proxy.
  - Go modül yapısının düzeltilmesi (`Go-CRM`).

- **Backend (API)**
  - `users`, `customers` ve `contacts` tabloları için veritabanı şemaları.
  - Müşteri (Customer) listeleme ve ekleme endpoint'leri.
  - İletişim (Contact) listeleme ve ekleme endpoint'leri.
  - Basit JWT tabanlı kullanıcı girişi (`/api/login`).

- **Frontend**
  - Vanilla JS ile kullanıcı giriş ekranı.
  - Müşteri ve iletişim kayıtlarını listeleyen, sekmeli dashboard arayüzü.
  - Müşteri ekleme ve düzenleme formu (modal).
  - Müşterilere iletişim kaydı ekleme formu (modal).
  - JavaScript modül yapısının düzeltilmesi (`type="module"`).
  
- **Test**
  - Unit ve entegrasyon testleri için temel yapı (`/tests` klasörü).
  - `Makefile` üzerinden testleri çalıştırmak için komutlar (`test-unit`, `test-integration`).

## Sırada Geliştirilecekler

- **Kimlik Doğrulama (Authentication)**
  - Gateway'de devre dışı bırakılan JWT/JWKS doğrulamasının düzgün bir şekilde yeniden etkinleştirilmesi.
  - Keycloak servisinin `docker-compose.yml`'e eklenmesi ve entegrasyonu.
  - Login işleminde kullanılan `demo123` gibi geçici şifre kontrolünün, `bcrypt` ile hash'lenmiş güvenli yapıya dönüştürülmesi.

- **API ve Fonksiyonellik**
  - Müşteri ve İletişimler için **Güncelleme (Update)** ve **Silme (Delete)** operasyonlarının tamamlanması.
  - API'ye daha detaylı hata yönetimi ve loglama eklenmesi.
  - API istekleri için "validation" (doğrulama) katmanı eklenmesi (örn: e-posta formatı kontrolü).

- **Testler**
  - Mevcut test yapısının, eklenen yeni özellikler (Customer ve Contact CRUD) için genişletilmesi.
  - Uçtan uca (E2E) testlerin yazılması.
  - Test kapsamını (coverage) artırma.

- **Kod Kalitesi ve Altyapı**
  - `go mod tidy` komutunun `git`'e ihtiyaç duymadan çalışmasının sağlanması veya `git`'in Docker imajına eklenmesi.
  - `Dockerfile`'ların daha verimli hale getirilmesi (multi-stage build optimizasyonu).

- **Observability**
  - Prometheus, Grafana, Jaeger gibi araçlarla metrik, log ve trace toplama altyapısının kurulması. 