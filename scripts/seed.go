package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Environment variable helper fonksiyonu
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Environment variables ile database konfigürasyonu
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "pass")
	dbName := getEnv("DB_NAME", "users")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Veritabanına bağlan
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB bağlantı hatası: %v", err)
	}
	defer db.Close()

	// Örnek kullanıcı bilgileri
	name := "Demo User"
	email := "demo@example.com"
	passwordHash := "demo123" // Gerçek projede bcrypt hash kullan!

	// Kullanıcıyı ekle (email unique olduğu için tekrar eklemez)
	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
	`, name, email, passwordHash)
	if err != nil {
		log.Fatalf("Kullanıcı eklenemedi: %v", err)
	}

	fmt.Println("Örnek kullanıcı başarıyla eklendi!")
}
