package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Veritabanına bağlan
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/users?sslmode=disable")
	if err != nil {
		log.Fatal("DB bağlantı hatası:", err)
	}
	defer db.Close()

	// Örnek kullanıcı bilgileri
	name := "Demo User"
	email := "demo@example.com"
	passwordHash := "$2a$10$demoHash" // bcrypt hash örneği

	// Kullanıcıyı ekle (email unique olduğu için tekrar eklemez)
	_, err = db.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
	`, name, email, passwordHash)
	if err != nil {
		log.Fatal("Kullanıcı eklenemedi:", err)
	}

	fmt.Println("Örnek kullanıcı başarıyla eklendi!")
}
