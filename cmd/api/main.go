package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"

	// Yeni customer handler importu
	"Go-CRM/pkg/common"
	"Go-CRM/pkg/customer"

	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	gootel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// İki veritabanı bağlantısını tutacak global değişkenler
var dbPrimary *sql.DB
var dbReplica *sql.DB

var jwtKey []byte

// Environment variable helper fonksiyonu
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Docker secrets dosyasından gizli bilgi okuma fonksiyonu
func readSecretFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Müşteri struct'ı
type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// İletişim struct'ı
type Contact struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

// JWT doğrulama middleware'i
func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Yetkilendirme gerekli", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Geçersiz veya süresi dolmuş token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Request-id logging middleware'i
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		w.Header().Set("X-Request-Id", reqID)
		ctx := context.WithValue(r.Context(), "request-id", reqID)
		// Trace-id propagation (OpenTelemetry varsa)
		traceID := ""
		if span := oteltrace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}
		ctx = context.WithValue(ctx, "trace-id", traceID)
		common.Logger.WithFields(logrus.Fields{
			"request_id": reqID,
			"trace_id":   traceID,
			"method":     r.Method,
			"path":       r.URL.Path,
		}).Info("request")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// healthzHandler endpoint'i için Swaggo açıklaması
// @Summary Healthcheck
// @Description Servisin canlı olup olmadığını kontrol eder
// @Tags Genel
// @Success 200 {object} map[string]string "Servis canlı"
// @Router /healthz [get]
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func initTracer() (func(), error) {
	endpoint := os.Getenv("JAEGER_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:14268/api/traces"
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-crm-api"),
		)),
	)
	gootel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(context.Background()) }, nil
}

func main() {
	common.InitLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	var err error

	// JWT secret'ı Docker secrets dosyasından oku
	jwtSecret, err := readSecretFile("/run/secrets/jwt_secret")
	if err != nil {
		log.Fatal("JWT secret dosyası okunamadı: ", err)
	}
	jwtKey = []byte(jwtSecret)

	// Ortam değişkenlerinden veritabanı URL'lerini al
	primaryURL := os.Getenv("DB_PRIMARY_URL")
	replicaURL := os.Getenv("DB_REPLICA_URL")

	if primaryURL == "" || replicaURL == "" {
		log.Fatal("DB_PRIMARY_URL ve DB_REPLICA_URL ortam değişkenleri ayarlanmalı!")
	}

	// DB şifresini Docker secrets dosyasından oku (örnek kullanım, istersen DB URL'lerinde kullanabilirsin)
	_, err = readSecretFile("/run/secrets/db_password")
	if err != nil {
		log.Println("DB password dosyası okunamadı, ortam değişkeni kullanılacak: ", err)
	}
	// DB URL'lerinde şifreyi dinamik olarak değiştirmek istersen:
	// primaryURL = strings.Replace(primaryURL, "pass", dbPassword, 1)
	// replicaURL = strings.Replace(replicaURL, "pass", dbPassword, 1)

	// Primary veritabanına bağlan (Yazma işlemleri için)
	dbPrimary, err = sql.Open("postgres", primaryURL)
	if err != nil {
		log.Fatalf("Primary veritabanına bağlanılamadı: %v", err)
	}
	defer dbPrimary.Close()

	// Replica veritabanına bağlan (Okuma işlemleri için)
	dbReplica, err = sql.Open("postgres", replicaURL)
	if err != nil {
		log.Fatalf("Replica veritabanına bağlanılamadı: %v", err)
	}
	defer dbReplica.Close()

	runMigrationsAndSeed(dbPrimary)

	// Redis bağlantısı başlatılıyor
	if err := common.InitRedis(); err != nil {
		log.Fatalf("Redis bağlantısı kurulamadı: %v", err)
	}
	log.Println("Redis bağlantısı başarılı.")
	common.InitRateLimiter(common.GetRedisClient())

	// Kafka bağlantısı başlatılıyor
	if err := common.InitKafka(); err != nil {
		log.Fatalf("Kafka bağlantısı kurulamadı: %v", err)
	}
	log.Println("Kafka bağlantısı başarılı.")

	shutdown, err := initTracer()
	if err != nil {
		log.Fatalf("Jaeger tracer başlatılamadı: %v", err)
	}
	defer shutdown()

	// Handler struct'ı oluşturuluyor
	handler := &customer.Handler{
		DBPrimary: dbPrimary,
		DBReplica: dbReplica,
	}

	router := mux.NewRouter()
	router.Use(requestIDMiddleware)

	// Swagger UI endpoint'i (JWT korumasız, public)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API rotaları
	router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		common.RateLimitMiddleware(http.HandlerFunc(handleLogin), 5, time.Minute).ServeHTTP(w, r)
	}).Methods("POST")
	router.HandleFunc("/healthz", healthzHandler).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// JWT korumalı alt router
	api := router.PathPrefix("/api").Subrouter()
	api.Use(jwtAuthMiddleware)

	// Müşteri işlemleri (yeni handler fonksiyonları)
	api.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			common.RateLimitMiddleware(http.HandlerFunc(handler.CreateCustomerHandler), 5, time.Minute).ServeHTTP(w, r)
			return
		}
		handler.GetCustomersHandler(w, r)
	}).Methods("GET", "POST")

	// İletişim kayıtları işlemleri (yeni handler fonksiyonları)
	api.HandleFunc("/contacts/{customerId}", handler.GetContactsHandler).Methods("GET")
	api.HandleFunc("/contacts", handler.CreateContactHandler).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-Id"},
		AllowCredentials: true,
	})
	log.Println("API sunucusu 8080 portunda başlatılıyor...")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}

// handleLogin endpoint'i için Swaggo açıklaması
// @Summary Kullanıcı girişi
// @Description Kullanıcı e-posta ve şifresiyle giriş yapar, JWT token döner
// @Tags Kimlik Doğrulama
// @Accept json
// @Produce json
// @Param login body User true "Giriş bilgileri"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {string} string "Geçersiz istek gövdesi"
// @Failure 401 {string} string "Kullanıcı adı veya şifre hatalı"
// @Router /api/login [post]
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Geçersiz istek gövdesi", http.StatusBadRequest)
		return
	}

	var storedUser User
	// Login işlemi kritik bir okuma olduğu için Primary'den yapılabilir,
	// veya anlık replikasyon gecikmesini kabul edip Replica'dan da yapılabilir.
	// Güvenilirlik için Primary'den okuyoruz.
	err := dbPrimary.QueryRow("SELECT id, password FROM users WHERE email=$1", creds.Email).Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Kullanıcı adı veya şifre hatalı", http.StatusUnauthorized)
		} else {
			http.Error(w, "Sunucu hatası", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Kullanıcı adı veya şifre hatalı", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: creds.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Müşteri işlemleri için handler fonksiyonları customer.Handler struct'ında tanımlı, örnek açıklamalar:
// GetCustomersHandler için:
// @Summary Müşteri listesini getir
// @Description Tüm müşterileri listeler
// @Tags Müşteri
// @Produce json
// @Success 200 {array} Customer
// @Failure 401 {string} string "Yetkilendirme gerekli"
// @Router /api/customers [get]

// CreateCustomerHandler için:
// @Summary Yeni müşteri oluştur
// @Description Yeni müşteri kaydı ekler
// @Tags Müşteri
// @Accept json
// @Produce json
// @Param customer body Customer true "Müşteri bilgileri"
// @Success 201 {object} Customer
// @Failure 400 {string} string "Geçersiz istek"
// @Failure 401 {string} string "Yetkilendirme gerekli"
// @Router /api/customers [post]

// GetContactsHandler için:
// @Summary Müşteriye ait iletişim kayıtlarını getir
// @Description Belirli bir müşterinin iletişim kayıtlarını listeler
// @Tags İletişim
// @Produce json
// @Param customerId path int true "Müşteri ID"
// @Success 200 {array} Contact
// @Failure 401 {string} string "Yetkilendirme gerekli"
// @Router /api/contacts/{customerId} [get]

// CreateContactHandler için:
// @Summary Yeni iletişim kaydı oluştur
// @Description Müşteriye yeni iletişim kaydı ekler
// @Tags İletişim
// @Accept json
// @Produce json
// @Param contact body Contact true "İletişim bilgileri"
// @Success 201 {object} Contact
// @Failure 400 {string} string "Geçersiz istek"
// @Failure 401 {string} string "Yetkilendirme gerekli"
// @Router /api/contacts [post]
// ... mevcut kod ...

// docker-init.sh script'i ve seeder.go'nun görevini üstlenen fonksiyon
// Sadece uygulama ilk ayağa kalktığında çalıştırılmalı.
func runMigrationsAndSeed(db *sql.DB) {
	log.Println("Migration'lar ve seeder kontrol ediliyor...")

	// 1. Tabloların var olup olmadığını kontrol et
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
	if err != nil {
		log.Fatalf("Tablo varlığı kontrol edilemedi: %v", err)
	}

	// Eğer 'users' tablosu zaten varsa, migration ve seed yapıldığını varsay ve çık.
	if exists {
		log.Println("Tablolar zaten mevcut. Migration ve seed atlanıyor.")
		return
	}

	log.Println("'users' tablosu bulunamadı. Migration'lar çalıştırılıyor...")

	// 2. Migration'ları çalıştır
	migrations := []string{
		`CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		);`,
		`CREATE TABLE customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			phone VARCHAR(50)
		);`,
		`CREATE TABLE contacts (
			id SERIAL PRIMARY KEY,
			customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			log.Fatalf("Migration başarısız: %v", err)
		}
	}
	log.Println("Tüm migration'lar başarıyla tamamlandı.")

	// 3. Seeder'ı çalıştır
	log.Println("Başlangıç verisi (seed) ekleniyor...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Şifre hash'lenemedi: %v", err)
	}

	email := "demo@example.com"
	_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, string(hashedPassword))
	if err != nil {
		log.Fatalf("Örnek kullanıcı eklenemedi: %v", err)
	}
	log.Println("Örnek kullanıcı başarıyla eklendi.")
}
