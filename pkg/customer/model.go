package customer

// Customer veri modeli
// Veritabanı ve API için ortak kullanılacak

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Contact struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
