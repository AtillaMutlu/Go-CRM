# Go Projesi Test Kılavuzu

Bu dokümantasyon, projenin test stratejisini ve testlerin nasıl çalıştırılacağını açıklar.

## Hızlı Başlangıç

Projedeki testler, host makineniz (bilgisayarınız) üzerinde çalıştırılacak şekilde tasarlanmıştır. Ancak, entegrasyon testleri gibi bazı testler, çalışan bir veritabanına ihtiyaç duyar.

Bu yüzden testleri çalıştırmadan önce **mutlaka** veritabanı servisini Docker ile ayağa kaldırmanız gerekir.

### Testleri Çalıştırma Adımları

**1. Veritabanını Başlatın**
Tüm sistemi (veritabanı dahil) başlatmak için proje ana dizininde aşağıdaki komutu çalıştırın:
```bash
make docker-up
```
Bu komut, entegrasyon testlerinin bağlanacağı PostgreSQL veritabanını başlatacaktır. Servislerin tamamen başlamasını bekleyin.

**2. Testleri Çalıştırın**
Veritabanı çalışırken, **yeni bir terminal açın** ve aşağıdaki `make` komutlarından istediğinizi çalıştırın:

- **Tüm Testler:**
  ```bash
  make test-all
  ```
  *Not: Bu komut bir `.sh` scripti çalıştırdığı için Windows'ta Git Bash veya WSL gerektirebilir.*

- **Sadece Unit Testler (Veritabanı GEREKTİRMEZ):**
  ```bash
  make test-unit
  ```

- **Sadece Entegrasyon Testleri (Veritabanı GEREKTİRİR):**
  ```bash
  make test-integration
  ```

- **Sadece E2E Testleri (Veritabanı GEREKTİRİR):**
  ```bash
  make test-e2e
  ```

## Test Mimarisi

Projemizde 3 temel test katmanı bulunmaktadır:

- **Unit Testleri:** Dış bağımlılığı olmayan, tekil fonksiyonları test eder. Hızlıdırlar ve sıkça çalıştırılmalıdır. (bkz: `tests/unit/`)
- **Entegrasyon Testleri:** Servislerin birbiriyle (özellikle API ve veritabanı) doğru şekilde entegre olup olmadığını test eder. (bkz: `tests/integration/`)
- **End-to-End (E2E) Testleri:** Tam bir kullanıcı senaryosunu baştan sona test eder. (bkz: `tests/e2e/`)

## Sorun Giderme

**Hata: `connection refused`**
- **Çözüm:** `make docker-up` komutunu çalıştırdığınızdan ve veritabanı konteynerinin sağlıklı bir şekilde başladığından emin olun. `docker ps` komutu ile `crm-postgres` konteynerinin "healthy" durumda olduğunu kontrol edebilirsiniz.

**Hata: `relation "..." does not exist`**
- **Çözüm:** Veritabanınızda migration'lar eksik olabilir. `make docker-down` ve ardından `make docker-up` komutlarını çalıştırarak konteynerlerin en baştan temiz bir şekilde kurulmasını ve migration'ların otomatik olarak çalışmasını sağlayın. 