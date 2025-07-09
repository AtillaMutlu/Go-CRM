# Go Mikroservis Monorepo

Bu repo, üretim-hazır Go tabanlı mikroservis iskeleti sunar. (CRM, stok yönetimi vb. için )

## Klasör Yapısı

- `cmd/`           : Her mikroservis için ana uygulama giriş noktası
- `pkg/`           : Ortak kütüphaneler ve handler fonksiyonları (örn: pkg/gateway)
- `deploy/`        : Dağıtım ve K8s manifest dosyaları
- `charts/`        : Helm chart'ları
- `configs/`       : Konfigürasyon dosyaları
- `proto/`         : Protobuf tanımları
- `migrations/`    : DB migration dosyaları (SQL up/down)
- `scripts/`       : Yardımcı scriptler (örn: seeder)
- `docs/`          : Dokümantasyon
- `tests/`         : Testler (unit, integration)
- `build/`         : Derlenmiş binary'ler
- `public/`        : Basit Bootstrap + JS frontend

## Handler ve Test Yapısı

- Tüm HTTP handler fonksiyonları `pkg/` altında modül bazlı tutulur (örn: `pkg/gateway/handlers.go`).
- Testler doğrudan bu fonksiyonları import ederek çalışır.

## Migration Yönetimi (Best Practice)

- Migration dosyalarını `migrations/` klasöründe, sıralı ve up/down SQL olarak tut.
- Migration işlemleri için [golang-migrate](https://github.com/golang-migrate/migrate) CLI kullanılır.
- Migration'ı uygulamak için:

```sh
make migrate-up
```
- Geri almak için:
```sh
make migrate-down
```
- Migration'lar kodla versionlanır, takımda herkes aynı şemayı kullanır.
- Migration'ı uygulama başında veya CI/CD pipeline'ında otomatik çalıştırabilirsin.

## Frontend + Backend Entegrasyonu (Örnek Mini CRM)

- `public/` altında Bootstrap + vanilla JS ile login, dashboard, müşteri CRUD ve iletişim listesi arayüzü hazır.
- `cmd/api/main.go` altında REST API (login, müşteri CRUD, iletişim) gerçek PostgreSQL ile çalışır.
- JWT ile login, tüm işlemler güvenli.
- Migration ve seeder ile tablo ve örnek kullanıcı otomatik oluşturulur.

### Hızlı Başlangıç (Tüm Akış)

```sh
# 1. PostgreSQL başlat
make migrate-up
# veya
# docker compose -f deploy/docker-compose.postgres.yml up

# 2. Migration ve seeder
make migrate-up
# Seeder ile örnek kullanıcı ekle
# (email: demo@example.com, şifre: demo123)
go run scripts/seed.go

# 3. API'yi başlat
make run SERVICE=api
# veya
go run cmd/api/main.go

# 4. Frontend'i aç (public/index.html)
# Login ol, müşteri CRUD ve iletişim işlemlerini test et
```

## Gereksinimler
- Go 1.22+
- Docker 25+
- [golang-migrate](https://github.com/golang-migrate/migrate) (migration için)

## Katkı
PR ve issue açabilirsiniz. Her modül için ADR (docs/adr) ekleyin. 