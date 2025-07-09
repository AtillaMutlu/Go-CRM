package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

var jwtSecret = []byte("supersecret") // Gerçek projede env'den al

var db *sql.DB

func main() {
	var err error
	// PostgreSQL bağlantısı
	db, err = sql.Open("postgres", "postgres://user:pass@localhost:5432/users?sslmode=disable")
	if err != nil {
		log.Fatal("DB bağlantı hatası:", err)
	}
	defer db.Close()

	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/customers", jwtAuth(customersHandler))
	http.HandleFunc("/api/customers/", jwtAuth(customerDetailHandler))
	http.HandleFunc("/api/contacts", jwtAuth(contactsHandler))

	fmt.Println("API servis 8085 portunda başlatıldı...")
	log.Fatal(http.ListenAndServe(":8085", nil))
}

// JWT doğrulama middleware
func jwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "Yetkisiz: Bearer token gerekli", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Yetkisiz: JWT geçersiz", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// /api/login (POST)
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}
	var id int
	var email, passwordHash string
	err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", req.Email).Scan(&id, &email, &passwordHash)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}
	// Demo için şifre hash kontrolü yok, gerçek projede bcrypt kullan!
	if req.Password != "demo123" && req.Password != passwordHash {
		http.Error(w, "Şifre hatalı", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(jwtSecret)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

// Müşteri struct'ı
type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// /api/customers (GET, POST)
func customersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rows, err := db.Query("SELECT id, name, email, COALESCE(phone, '') FROM customers ORDER BY id DESC")
		if err != nil {
			http.Error(w, "DB hatası", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var customers []Customer
		for rows.Next() {
			var c Customer
			if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone); err == nil {
				customers = append(customers, c)
			}
		}
		json.NewEncoder(w).Encode(customers)
		return
	}
	if r.Method == http.MethodPost {
		var c Customer
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Geçersiz veri", http.StatusBadRequest)
			return
		}
		var id int
		err := db.QueryRow("INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING id", c.Name, c.Email, c.Phone).Scan(&id)
		if err != nil {
			http.Error(w, "DB ekleme hatası", http.StatusInternalServerError)
			return
		}
		c.ID = id
		json.NewEncoder(w).Encode(c)
		return
	}
	http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
}

// /api/customers/{id} (PUT, DELETE)
func customerDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/customers/")
	if id == "" {
		http.Error(w, "ID gerekli", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodPut {
		var c Customer
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Geçersiz veri", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("UPDATE customers SET name=$1, email=$2, phone=$3 WHERE id=$4", c.Name, c.Email, c.Phone, id)
		if err != nil {
			http.Error(w, "DB güncelleme hatası", http.StatusInternalServerError)
			return
		}
		c.ID = atoi(id)
		json.NewEncoder(w).Encode(c)
		return
	}
	if r.Method == http.MethodDelete {
		_, err := db.Exec("DELETE FROM customers WHERE id=$1", id)
		if err != nil {
			http.Error(w, "DB silme hatası", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
}

// İletişim struct'ı
type Contact struct {
	ID           int    `json:"id"`
	CustomerName string `json:"customer_name"`
	Message      string `json:"message"`
	Date         string `json:"date"`
}

// /api/contacts (GET)
func contactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query(`SELECT c.id, cu.name, c.message, to_char(c.created_at, 'YYYY-MM-DD') FROM contacts c JOIN customers cu ON c.customer_id = cu.id ORDER BY c.id DESC`)
	if err != nil {
		http.Error(w, "DB hatası", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var contacts []Contact
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.ID, &c.CustomerName, &c.Message, &c.Date); err == nil {
			contacts = append(contacts, c)
		}
	}
	json.NewEncoder(w).Encode(contacts)
}

func atoi(s string) int {
	n, _ := fmt.Sscanf(s, "%d", new(int))
	return n
}
