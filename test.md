# Go Mikroservis Test Kılavuzu

Bu dokümantasyon, Go mikroservis projesinin kapsamlı test stratejisini ve implementasyonunu açıklar.

## 📋 İçindekiler

1. [Hızlı Başlangıç](#hızlı-başlangıç)
2. [Test Mimarisi](#test-mimarisi)
3. [Test Katmanları](#test-katmanları)
4. [Test Komutları](#test-komutları)
5. [Test Coverage](#test-coverage)
6. [CI/CD Entegrasyonu](#cicd-entegrasyonu)
7. [Troubleshooting](#troubleshooting)
8. [Performance Testing](#performance-testing)
9. [Security Testing](#security-testing)

## 🚀 Hızlı Başlangıç

### Ön Koşullar

```bash
# PostgreSQL başlat
docker compose -f deploy/docker-compose.postgres.yml up -d

# Migration tool'u kur (eğer yoksa)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Veritabanını hazırla
make migrate-up
go run scripts/seed.go
```

### Tüm Testleri Çalıştır

```bash
# Tek komutla tüm test suite'ini çalıştır
make test-all
```

### Adım Adım Test Etme

```bash
make test-infrastructure  # Infrastructure kontrolü
make test-unit           # Unit tests
make test-integration    # Database testleri
make test-e2e           # End-to-end tests
```

## 🏗️ Test Mimarisi

Projemizde **4 katmanlı test piramidi** stratejisi kullanılmaktadır:

```
        /\
       /  \
      /E2E \     ← End-to-End Tests (Az sayıda, yavaş)
     /______\
    /        \
   /Integration\ ← Integration Tests (Orta seviye)
  /____________\
 /              \
/   Unit Tests   \ ← Unit Tests (Çok sayıda, hızlı)
\________________/
\  Infrastructure / ← Infrastructure Tests (Temel)
 \______________/
```

### Test Dosya Yapısı

```
tests/
├── unit/                 # Unit testler
│   └── gateway_test.go
├── integration/          # Integration testler
│   └── api_test.go
├── e2e/                 # End-to-end testler
│   └── user_flow_test.go
└── ...

scripts/
├── test-infrastructure.sh  # Infrastructure test
├── run-all-tests.sh        # Master test runner
└── ...

migrations/
└── test/                   # Test-specific migrations
    └── 002_add_customers_table.up.sql
```

## 🧪 Test Katmanları

### 1. Infrastructure Tests

**Amaç:** Temel altyapının çalışır durumda olduğunu doğrula

**Dosya:** `scripts/test-infrastructure.sh`

**Test eder:**
- ✅ PostgreSQL bağlantısı
- ✅ Migration sisteminin çalışması
- ✅ Seed data oluşturulması  
- ✅ Gateway health endpoint
- ✅ API login endpoint
- ✅ Temel servis ayakta olma durumu

**Çalıştırma:**
```bash
make test-infrastructure
# veya
bash scripts/test-infrastructure.sh
```

**Beklenen Çıktı:**
```
🧪 Infrastructure Test Başlatılıyor...
📊 PostgreSQL: ✅ BAŞARILI
📋 Migration: ✅ BAŞARILI
🌱 Seed Data: ✅ BAŞARILI
🚪 Gateway Health: ✅ BAŞARILI
📡 API Login: ✅ BAŞARILI
🎉 Tüm infrastructure testleri BAŞARILI!
```

### 2. Unit Tests

**Amaç:** İzole fonksiyonları ve handler'ları test et

**Dosya:** `tests/unit/gateway_test.go`

**Test eder:**
- ✅ HTTP handler fonksiyonları
- ✅ Middleware logic
- ✅ Rate limiting simulation
- ✅ Concurrent access patterns
- ✅ Performance benchmarks

**Özellikler:**
- HTTP test server kullanımı
- Goroutine-safe testler
- Benchmark testleri
- Error handling testleri

**Çalıştırma:**
```bash
make test-unit
# veya
go test -v ./tests/unit/...
```

**Örnek Test Çıktısı:**
```
=== RUN   TestHealthzHandler
--- PASS: TestHealthzHandler (0.00s)
=== RUN   TestHealthzConcurrency
--- PASS: TestHealthzConcurrency (0.01s)
PASS
ok      tests/unit    0.123s
```

### 3. Integration Tests

**Amaç:** Database ve API katmanları entegrasyonunu test et

**Dosya:** `tests/integration/api_test.go`

**Test eder:**
- ✅ Database CRUD işlemleri
- ✅ Transaction rollback
- ✅ Login entegrasyonu
- ✅ Database connection stability
- ✅ Data integrity

**Özellikler:**
- Gerçek PostgreSQL bağlantısı
- Test-specific database kullanımı
- Automatic cleanup
- Transaction testing

**Çalıştırma:**
```bash
make test-integration
# veya
export TEST_DB_URL="postgres://user:pass@localhost:5432/users?sslmode=disable"
go test -v ./tests/integration/...
```

### 4. End-to-End Tests

**Amaç:** Tam kullanıcı akışını simüle et

**Dosya:** `tests/e2e/user_flow_test.go`

**Test eder:**
- ✅ Complete authentication flow
- ✅ Customer CRUD operations
- ✅ Authorization checks
- ✅ Error handling flows
- ✅ Basic load testing
- ✅ Security validations

**Test Senaryoları:**

1. **Authentication Flow**
   - Valid login → JWT token alma
   - Invalid login → Error handling

2. **Customer Management**
   - Create customer
   - List customers
   - Update customer
   - Delete customer

3. **Authorization Tests**
   - Unauthorized access attempts
   - Invalid token handling

4. **Error Handling**
   - Invalid JSON requests
   - Non-existent resource operations

5. **Performance Tests**
   - Concurrent request handling

**Çalıştırma:**
```bash
make test-e2e
# veya
go test -v ./tests/e2e/...
```

## 📋 Test Komutları

### Temel Komutlar

| Komut | Açıklama | Süre | Dependencies |
|-------|----------|------|--------------|
| `make test-infrastructure` | Infrastructure kontrolü | ~30s | PostgreSQL |
| `make test-unit` | Unit testler | ~5s | - |
| `make test-integration` | Database testleri | ~10s | PostgreSQL + Migrations |
| `make test-e2e` | End-to-end akış | ~15s | API servis + Database |
| `make test-all` | Tüm testler sırasıyla | ~60s | Tümü |

### İleri Seviye Komutlar

```bash
# Performance benchmarks
make test-performance
go test -bench=. -benchmem ./tests/unit/...

# Test coverage raporu
make test-coverage
# Çıktı: coverage.html dosyası

# Specific test çalıştırma
go test -v -run TestSpecificTest ./tests/unit/...

# Parallel testing
go test -parallel 4 ./tests/...

# Test cache temizleme
go clean -testcache

# Verbose output ile debugging
go test -v -race ./tests/...
```

### Test Migration Komutları

```bash
# Test-specific migrations
make migrate-test-up    # Test tabloları oluştur
make migrate-test-down  # Test tabloları kaldır

# Manuel migration
migrate -path migrations/test -database "postgres://user:pass@localhost:5432/users?sslmode=disable" up
```

## 📊 Test Coverage

### Coverage Raporu Oluşturma

```bash
# HTML coverage raporu
make test-coverage

# Terminal'de coverage
go test -cover ./...

# Detaylı coverage profili
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Coverage Hedefleri

| Test Tipi | Coverage Hedefi | Mevcut |
|-----------|----------------|---------|
| Unit Tests | >80% | TBD |
| Integration | >70% | TBD |
| E2E | >60% | TBD |
| Overall | >75% | TBD |

### Coverage Analizi

```bash
# Package bazlı coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Function bazlı coverage
go tool cover -func=coverage.out

# Coverage'ı JSON formatında export
go tool cover -json=coverage.out
```

## 🔄 CI/CD Entegrasyonu

### GitHub Actions Örneği

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: pass
          POSTGRES_USER: user
          POSTGRES_DB: users
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.22
        
    - name: Install migrate tool
      run: |
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        
    - name: Run Infrastructure Tests
      run: make test-infrastructure
      
    - name: Run Unit Tests
      run: make test-unit
      
    - name: Run Integration Tests
      run: make test-integration
      env:
        TEST_DB_URL: postgres://user:pass@localhost:5432/users?sslmode=disable
        
    - name: Run E2E Tests
      run: make test-e2e
      
    - name: Generate Coverage Report
      run: make test-coverage
      
    - name: Upload Coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### Docker-based Testing

```bash
# Test container oluştur
docker build -f Dockerfile.test -t app-test .

# Container içinde testleri çalıştır
docker run --rm -v $(pwd):/app app-test make test-all

# Docker Compose ile test environment
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 🔧 Troubleshooting

### Yaygın Sorunlar ve Çözümleri

#### 1. PostgreSQL Bağlantı Hatası

**Hata:**
```
connection refused to localhost:5432
```

**Çözüm:**
```bash
# PostgreSQL'i başlat
docker compose -f deploy/docker-compose.postgres.yml up -d

# Bağlantıyı test et
psql -h localhost -U user -d users -c "SELECT 1;"

# Port kontrolü
netstat -tulpn | grep 5432
```

#### 2. Migration Tool Hatası

**Hata:**
```
migrate: command not found
```

**Çözüm:**
```bash
# Migration tool'u kur
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# PATH kontrolü
export PATH=$PATH:$(go env GOPATH)/bin
```

#### 3. Test Database Hatası

**Hata:**
```
relation "customers" does not exist
```

**Çözüm:**
```bash
# Migration'ları çalıştır
make migrate-up
make migrate-test-up

# Seed data ekle
go run scripts/seed.go
```

#### 4. JWT Token Hatası

**Hata:**
```
JWT invalid
```

**Çözüm:**
```bash
# Seed user'ı kontrol et
psql -h localhost -U user -d users -c "SELECT * FROM users;"

# Demo user ekle
go run scripts/seed.go

# Login endpoint test et
curl -X POST http://localhost:8085/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"demo123"}'
```

#### 5. Port Conflict Hatası

**Hata:**
```
bind: address already in use
```

**Çözüm:**
```bash
# Port kullanımını kontrol et
lsof -i :8085
lsof -i :8080

# Process'i durdur
kill -9 $(lsof -t -i:8085)

# Alternative port kullan
PORT=8086 go run cmd/api/main.go
```

### Debug Komutları

```bash
# Detaylı test output
go test -v -race ./tests/...

# Specific test debug
go test -v -run TestCustomerCRUD ./tests/integration/...

# Test timing
go test -v -timeout 30s ./tests/...

# Memory profiling
go test -memprofile=mem.prof ./tests/...
go tool pprof mem.prof

# CPU profiling
go test -cpuprofile=cpu.prof ./tests/...
go tool pprof cpu.prof
```

### Log Analizi

```bash
# Test loglarını yakala
go test -v ./tests/... 2>&1 | tee test.log

# Error'ları filtrele
grep -i "error\|fail\|fatal" test.log

# Test timing analizi
grep -E "PASS|FAIL|RUN" test.log
```

## 🚀 Performance Testing

### Benchmark Tests

```bash
# Tüm benchmark'ları çalıştır
make test-performance

# Specific benchmark
go test -bench=BenchmarkHealthzHandler ./tests/unit/...

# Memory allocation profiling
go test -bench=. -benchmem ./tests/unit/...

# CPU profiling ile benchmark
go test -bench=. -cpuprofile=cpu.prof ./tests/unit/...
```

### Load Testing

#### Apache Bench ile

```bash
# Basic load test
ab -n 1000 -c 10 http://localhost:8085/api/customers

# Authentication ile
ab -n 100 -c 5 -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8085/api/customers

# POST request test
ab -n 100 -c 5 -p customer.json -T application/json \
  http://localhost:8085/api/customers
```

#### Custom Load Test

```bash
# scripts/load-test.sh
#!/bin/bash
echo "Load testing başlatılıyor..."

# API'yi başlat
go run cmd/api/main.go &
API_PID=$!
sleep 3

# Progressively increase load
for concurrent in 1 5 10 20 50; do
  echo "Testing with $concurrent concurrent users..."
  ab -n $((concurrent * 10)) -c $concurrent http://localhost:8085/api/login \
    -p login.json -T application/json
  sleep 2
done

kill $API_PID
```

### Performance Metrics

```bash
# Response time monitoring
curl -o /dev/null -s -w "Time: %{time_total}s\n" \
  http://localhost:8085/api/customers

# Throughput testing
wrk -t12 -c400 -d30s http://localhost:8085/api/customers

# Memory usage monitoring
while true; do
  ps aux | grep "go run cmd/api/main.go"
  sleep 1
done
```

## 🔒 Security Testing

### Authentication Tests

```bash
# JWT bypass attempt
curl -H "Authorization: Bearer fake.jwt.token" \
  http://localhost:8085/api/customers

# No token test
curl http://localhost:8085/api/customers

# Malformed token
curl -H "Authorization: InvalidTokenFormat" \
  http://localhost:8085/api/customers
```

### SQL Injection Tests

```bash
# SQL injection in login
curl -X POST http://localhost:8085/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com'\'' OR 1=1--","password":"test"}'

# SQL injection in customer creation
curl -X POST http://localhost:8085/api/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test'\''DROP TABLE customers--","email":"test@test.com"}'
```

### Input Validation Tests

```bash
# XSS attempt
curl -X POST http://localhost:8085/api/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"<script>alert(\"XSS\")</script>","email":"test@test.com"}'

# Oversized payload
curl -X POST http://localhost:8085/api/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"'$(yes A | head -n 10000 | tr -d '\n')'","email":"test@test.com"}'
```

### Rate Limiting Tests

```bash
# Rate limit test script
for i in {1..15}; do
  curl -o /dev/null -s -w "%{http_code}\n" \
    http://localhost:8080/api
done
# Beklenen: İlk 10 istek 200, sonraki 5 istek 429
```

## 📈 Test Metrics & Reporting

### Test Execution Metrics

```bash
# Test execution time tracking
time make test-all

# Individual test timing
go test -v ./tests/unit/... | grep "PASS\|FAIL" | awk '{print $3, $2}'

# Test success rate tracking
go test ./tests/... 2>&1 | grep -E "PASS|FAIL" | sort | uniq -c
```

### Test Result Formats

```bash
# JSON format test results
go test -json ./tests/... > test-results.json

# JUnit XML format (third-party tool gerekli)
go get -u github.com/jstemmer/go-junit-report
go test ./tests/... -v 2>&1 | go-junit-report > test-results.xml

# Test summary raporu
go test ./tests/... -v | grep -E "PASS|FAIL|RUN" | \
  awk '/^=== RUN/ {test=$3} /^--- PASS/ {pass++} /^--- FAIL/ {fail++} 
       END {print "Total:", pass+fail, "Passed:", pass, "Failed:", fail}'
```

### Monitoring & Alerting

```bash
# Test failure alerting script
#!/bin/bash
if ! make test-all; then
  echo "❌ Tests failed!" | mail -s "Test Failure Alert" admin@company.com
  # Slack webhook
  curl -X POST -H 'Content-type: application/json' \
    --data '{"text":"🚨 Test suite failed on commit '$(git rev-parse --short HEAD)'"}' \
    $SLACK_WEBHOOK_URL
fi
```

## 🎯 Best Practices

### Test Yazma Standartları

1. **AAA Pattern (Arrange, Act, Assert)**
```go
func TestCustomerCreation(t *testing.T) {
    // Arrange
    customer := Customer{Name: "Test", Email: "test@test.com"}
    
    // Act
    result, err := CreateCustomer(customer)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
}
```

2. **Table-Driven Tests**
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {"valid email", "test@test.com", true},
        {"invalid email", "invalid", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsValidEmail(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

3. **Test Cleanup**
```go
func TestWithCleanup(t *testing.T) {
    // Setup
    db := setupTestDB()
    
    // Cleanup
    t.Cleanup(func() {
        db.Close()
        cleanupTestData()
    })
    
    // Test logic...
}
```

### Test Environment Management

```bash
# Environment variables for tests
export TEST_ENV=true
export TEST_DB_URL="postgres://user:pass@localhost:5432/test_db"
export LOG_LEVEL=debug

# Test-specific configuration
cp config.test.yaml config.yaml

# Isolated test runs
go test -count=1 ./tests/... # Disable test caching
```

### Continuous Improvement

```bash
# Test flakiness detection
for i in {1..10}; do
  echo "Run $i:"
  go test ./tests/integration/... || echo "FLAKY TEST DETECTED"
done

# Performance regression testing
go test -bench=. ./tests/... > benchmark-old.txt
# After changes:
go test -bench=. ./tests/... > benchmark-new.txt
benchcmp benchmark-old.txt benchmark-new.txt
```

## 📚 Ekler

### Test Data Management

**Test Customer Data:**
```json
{
  "name": "Test Müşteri",
  "email": "test@example.com",
  "phone": "+905551234567"
}
```

**Test Login Credentials:**
```json
{
  "email": "demo@example.com",
  "password": "demo123"
}
```

### Test Utilities

**HTTP Test Helper:**
```go
func makeTestRequest(method, url string, body interface{}, token string) (*http.Response, error) {
    // Implementation in tests/e2e/user_flow_test.go
}
```

**Database Test Helper:**
```go
func setupTestDB() (*sql.DB, error) {
    // Implementation in tests/integration/api_test.go
}
```

### Reference Links

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Framework](https://github.com/stretchr/testify)
- [PostgreSQL Test Database](https://www.postgresql.org/docs/current/regress.html)
- [golang-migrate](https://github.com/golang-migrate/migrate)

---

## 🎉 Sonuç

Bu test suite'i ile mikroservis projenizi **production-ready** kalitede test edebilirsiniz:

✅ **Infrastructure Stability** - Altyapı kontrolü
✅ **Code Quality** - Unit test coverage  
✅ **Data Integrity** - Database entegrasyonu
✅ **User Experience** - End-to-end akışlar
✅ **Performance** - Load testing ve benchmarks
✅ **Security** - Auth ve validation testleri

**Hızlı Başlangıç:** `make test-all`
**Development Cycle:** `make test-unit`
**CI/CD Pipeline:** Tüm test katmanları

Bu kapsamlı test stratejisi sayesinde mikroservis projenizin kalitesini garanti altına alabilir ve güvenle production'a deploy edebilirsiniz! 🚀 