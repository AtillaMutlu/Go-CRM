package customer

import (
	"database/sql"
)

// Yeni müşteri ekler
func createCustomerRepo(db *sql.DB, c *Customer) error {
	return db.QueryRow(
		"INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING id",
		c.Name, c.Email, c.Phone,
	).Scan(&c.ID)
}

// Belirli bir müşterinin iletişim kayıtlarını getirir (pagination)
func getContactsByCustomerIDRepo(db *sql.DB, params ContactListParams) (ContactListResult, error) {
	query := "SELECT id, customer_id, content, created_at FROM contacts WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	rows, err := db.Query(query, params.CustomerID, params.PageSize, (params.Page-1)*params.PageSize)
	if err != nil {
		return ContactListResult{}, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		var createdAt sql.NullTime
		if err := rows.Scan(&contact.ID, &contact.CustomerID, &contact.Content, &createdAt); err != nil {
			return ContactListResult{}, err
		}
		if createdAt.Valid {
			contact.CreatedAt = createdAt.Time.Format("02-01-2006 15:04")
		}
		contacts = append(contacts, contact)
	}

	// Toplam kayıt sayısı
	var total int
	countQ := "SELECT COUNT(*) FROM contacts WHERE customer_id = $1"
	if err := db.QueryRow(countQ, params.CustomerID).Scan(&total); err != nil {
		return ContactListResult{}, err
	}

	return ContactListResult{
		Contacts: contacts,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// Yeni iletişim kaydı ekler
func createContactRepo(db *sql.DB, contact *Contact) error {
	return db.QueryRow(
		"INSERT INTO contacts (customer_id, content) VALUES ($1, $2) RETURNING id, created_at",
		contact.CustomerID, contact.Content,
	).Scan(&contact.ID, &contact.CreatedAt)
}
