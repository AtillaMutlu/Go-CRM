package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

// Test veritabanı bağlantısı
func setupTestDB() (*sql.DB, error) {
	// Test için ayrı veritabanı kullan
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:pass@localhost:5432/users_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Test tabloları oluştur
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);
		
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			phone VARCHAR(20),
			created_at TIMESTAMP DEFAULT now()
		);
		
		CREATE TABLE IF NOT EXISTS contacts (
			id SERIAL PRIMARY KEY,
			customer_id INTEGER REFERENCES customers(id),
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);
	`)

	return db, err
}

// Test cleanup
func cleanupTestDB(db *sql.DB) {
	db.Exec("DELETE FROM contacts")
	db.Exec("DELETE FROM customers")
	db.Exec("DELETE FROM users")
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		fmt.Printf("Test DB setup hatası: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	// Test'leri çalıştır
	code := m.Run()

	// Cleanup
	cleanupTestDB(testDB)
	os.Exit(code)
}

// Customer CRUD testleri
func TestCustomerCRUD(t *testing.T) {
	// Test data
	customer := map[string]interface{}{
		"name":  "Test Müşteri",
		"email": "test@example.com",
		"phone": "+90123456789",
	}

	// CREATE test
	t.Run("CreateCustomer", func(t *testing.T) {
		jsonData, _ := json.Marshal(customer)
		req := httptest.NewRequest("POST", "/api/customers", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-test-token") // Mock JWT

		// Bu kısımda actual API handler'ını çağırabilirsin
		// Şimdilik DB'ye direkt yazıyoruz
		var id int
		err := testDB.QueryRow(`
			INSERT INTO customers (name, email, phone) 
			VALUES ($1, $2, $3) RETURNING id
		`, customer["name"], customer["email"], customer["phone"]).Scan(&id)

		if err != nil {
			t.Fatalf("Customer oluşturulamadı: %v", err)
		}

		if id == 0 {
			t.Fatal("Customer ID alınamadı")
		}

		t.Logf("✅ Customer oluşturuldu, ID: %d", id)
	})

	// READ test
	t.Run("GetCustomers", func(t *testing.T) {
		rows, err := testDB.Query("SELECT id, name, email, COALESCE(phone, '') FROM customers")
		if err != nil {
			t.Fatalf("Customer listesi alınamadı: %v", err)
		}
		defer rows.Close()

		customers := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var name, email, phone string
			if err := rows.Scan(&id, &name, &email, &phone); err == nil {
				customers = append(customers, map[string]interface{}{
					"id": id, "name": name, "email": email, "phone": phone,
				})
			}
		}

		if len(customers) == 0 {
			t.Fatal("Customer listesi boş")
		}

		t.Logf("✅ %d customer bulundu", len(customers))
	})

	// UPDATE test
	t.Run("UpdateCustomer", func(t *testing.T) {
		// Önce bir customer bul
		var customerID int
		err := testDB.QueryRow("SELECT id FROM customers LIMIT 1").Scan(&customerID)
		if err != nil {
			t.Skip("Update testi için customer bulunamadı")
		}

		updatedData := map[string]interface{}{
			"name":  "Güncellenmiş Müşteri",
			"email": "updated@example.com",
			"phone": "+90987654321",
		}

		_, err = testDB.Exec(`
			UPDATE customers SET name=$1, email=$2, phone=$3 WHERE id=$4
		`, updatedData["name"], updatedData["email"], updatedData["phone"], customerID)

		if err != nil {
			t.Fatalf("Customer güncellenemedi: %v", err)
		}

		// Güncellemeyi doğrula
		var name string
		err = testDB.QueryRow("SELECT name FROM customers WHERE id=$1", customerID).Scan(&name)
		if err != nil || name != updatedData["name"] {
			t.Fatalf("Customer güncelleme doğrulanamadı: %v", err)
		}

		t.Logf("✅ Customer güncellendi, ID: %d", customerID)
	})

	// DELETE test
	t.Run("DeleteCustomer", func(t *testing.T) {
		// Test customer oluştur
		var customerID int
		err := testDB.QueryRow(`
			INSERT INTO customers (name, email, phone) 
			VALUES ('Silinecek Müşteri', 'delete@example.com', '123') 
			RETURNING id
		`).Scan(&customerID)
		if err != nil {
			t.Fatalf("Test customer oluşturulamadı: %v", err)
		}

		// Customer'ı sil
		result, err := testDB.Exec("DELETE FROM customers WHERE id=$1", customerID)
		if err != nil {
			t.Fatalf("Customer silinemedi: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected != 1 {
			t.Fatalf("Beklenmeyen silme sonucu: %d satır etkilendi", rowsAffected)
		}

		t.Logf("✅ Customer silindi, ID: %d", customerID)
	})
}

// Login entegrasyonu test
func TestLoginIntegration(t *testing.T) {
	// Test user oluştur
	_, err := testDB.Exec(`
		INSERT INTO users (name, email, password_hash) 
		VALUES ('Test User', 'test@example.com', 'test123')
		ON CONFLICT (email) DO NOTHING
	`)
	if err != nil {
		t.Fatalf("Test user oluşturulamadı: %v", err)
	}

	// Login testi
	t.Run("ValidLogin", func(t *testing.T) {
		var email, passwordHash string
		err := testDB.QueryRow(`
			SELECT email, password_hash FROM users WHERE email=$1
		`, "test@example.com").Scan(&email, &passwordHash)

		if err != nil {
			t.Fatalf("User bulunamadı: %v", err)
		}

		if email != "test@example.com" {
			t.Fatalf("Yanlış email: %v", email)
		}

		t.Logf("✅ Login testi başarılı: %s", email)
	})

	t.Run("InvalidLogin", func(t *testing.T) {
		var count int
		err := testDB.QueryRow(`
			SELECT COUNT(*) FROM users WHERE email=$1
		`, "nonexistent@example.com").Scan(&count)

		if err != nil {
			t.Fatalf("DB sorgu hatası: %v", err)
		}

		if count != 0 {
			t.Fatal("Var olmayan user bulundu!")
		}

		t.Log("✅ Invalid login testi başarılı")
	})
}

// Database connection test
func TestDatabaseConnection(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Fatalf("Database bağlantısı başarısız: %v", err)
	}

	t.Log("✅ Database bağlantısı başarılı")
}

// Transaction test
func TestTransactionRollback(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("Transaction başlatılamadı: %v", err)
	}

	// Test data ekle
	_, err = tx.Exec(`
		INSERT INTO customers (name, email, phone) 
		VALUES ('Transaction Test', 'tx@example.com', '123')
	`)
	if err != nil {
		t.Fatalf("Transaction insert hatası: %v", err)
	}

	// Rollback yap
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("Transaction rollback hatası: %v", err)
	}

	// Data'nın eklenmediğini doğrula
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM customers WHERE email='tx@example.com'").Scan(&count)
	if err != nil {
		t.Fatalf("Rollback doğrulama hatası: %v", err)
	}

	if count != 0 {
		t.Fatal("Transaction rollback çalışmadı!")
	}

	t.Log("✅ Transaction rollback testi başarılı")
}
