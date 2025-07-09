package customer

import (
	"Go-CRM/pkg/common"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Pagination ve filtreleme için parametreler
type CustomerListParams struct {
	Page     int
	PageSize int
	Search   string // İsim veya e-posta araması
}

type CustomerListResult struct {
	Customers []Customer
	Total     int
	Page      int
	PageSize  int
}

// Müşteri listeleme (pagination + filtreleme)
func GetCustomers(db *sql.DB, params CustomerListParams) (CustomerListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("customers:%s:%d:%d", params.Search, params.Page, params.PageSize)
	if cached, err := common.RedisGet(ctx, cacheKey); err == nil && cached != "" {
		var result CustomerListResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	var (
		args   []interface{}
		where  []string
		query  = "SELECT id, name, email, phone FROM customers"
		countQ = "SELECT COUNT(*) FROM customers"
	)
	if params.Search != "" {
		where = append(where, "(LOWER(name) LIKE $1 OR LOWER(email) LIKE $1)")
		args = append(args, "%"+strings.ToLower(params.Search)+"%")
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
		countQ += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY id DESC LIMIT $2 OFFSET $3"
	args = append(args, params.PageSize, (params.Page-1)*params.PageSize)

	rows, err := db.Query(query, args...)
	if err != nil {
		return CustomerListResult{}, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone); err != nil {
			return CustomerListResult{}, err
		}
		customers = append(customers, c)
	}

	// Toplam kayıt sayısı
	total := 0
	if err := db.QueryRow(countQ, args[:len(args)-2]...).Scan(&total); err != nil {
		return CustomerListResult{}, err
	}

	result := CustomerListResult{
		Customers: customers,
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
	}

	// Sonucu cache'e yaz
	if b, err := json.Marshal(result); err == nil {
		_ = common.RedisSet(ctx, cacheKey, string(b), 30*time.Second)
	}

	return result, nil
}

// Müşteri ekleme (validasyon ve güvenlik)
func CreateCustomer(db *sql.DB, c *Customer) error {
	if err := validateCustomer(*c); err != nil {
		return err
	}
	if err := createCustomerRepo(db, c); err != nil {
		return err
	}
	// Kafka event publish
	go func(cust Customer) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		payload, _ := json.Marshal(cust)
		_ = common.PublishEvent(ctx, fmt.Sprintf("customer-%d", cust.ID), string(payload))
	}(*c)
	return nil
}

// İletişim kayıtlarını getir (pagination)
func GetContactsByCustomerID(db *sql.DB, params ContactListParams) (ContactListResult, error) {
	if params.CustomerID <= 0 {
		return ContactListResult{}, errors.New("Geçersiz müşteri ID")
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return getContactsByCustomerIDRepo(db, params)
}

// İletişim kaydı ekleme (validasyon)
func CreateContact(db *sql.DB, contact *Contact) error {
	if contact.CustomerID <= 0 || utf8.RuneCountInString(contact.Content) < 3 {
		return errors.New("Geçersiz iletişim kaydı")
	}
	return createContactRepo(db, contact)
}

// Pagination ve filtreleme için parametreler
type ContactListParams struct {
	CustomerID int
	Page       int
	PageSize   int
}

type ContactListResult struct {
	Contacts []Contact
	Total    int
	Page     int
	PageSize int
}

// --- Validasyon Fonksiyonları ---
func validateCustomer(c Customer) error {
	if utf8.RuneCountInString(c.Name) < 2 {
		return errors.New("İsim en az 2 karakter olmalı")
	}
	if c.Email != "" && !strings.Contains(c.Email, "@") {
		return errors.New("Geçersiz e-posta adresi")
	}
	if c.Phone != "" && len(c.Phone) < 7 {
		return errors.New("Telefon numarası çok kısa")
	}
	return nil
}

// Not: Customer ve Contact listeleme sorgularında sadece gerekli alanlar çekiliyor, gereksiz sütun sorgusu yoktur.
