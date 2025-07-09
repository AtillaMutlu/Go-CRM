package customer

import (
	"Go-CRM/pkg/common"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler fonksiyonları sade tutulur, iş mantığı service katmanında

type Handler struct {
	DBPrimary *sql.DB
	DBReplica *sql.DB
}

// Müşteri listeleme (GET /api/customers?page=1&pageSize=10&search=ali)
func (h *Handler) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	params := CustomerListParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}
	result, err := GetCustomers(h.DBReplica, params)
	if err != nil {
		http.Error(w, "Müşteri listesi alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Müşteri ekleme (POST /api/customers)
func (h *Handler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		common.WriteError(w, http.StatusBadRequest, "Geçersiz istek gövdesi", err)
		return
	}
	if !common.IsEmailValid(c.Email) {
		common.WriteError(w, http.StatusBadRequest, "Geçersiz e-posta adresi", nil)
		return
	}
	if !common.IsPhoneValid(c.Phone) {
		common.WriteError(w, http.StatusBadRequest, "Geçersiz telefon numarası", nil)
		return
	}
	if err := CreateCustomer(h.DBPrimary, &c); err != nil {
		common.WriteError(w, http.StatusBadRequest, "Müşteri eklenemedi", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// Belirli bir müşterinin iletişim kayıtları (GET /api/contacts/{customerId})
func (h *Handler) GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	customerIDStr := strings.TrimPrefix(r.URL.Path, "/api/contacts/")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil || customerID <= 0 {
		common.WriteError(w, http.StatusBadRequest, "Geçersiz müşteri ID", nil)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	params := ContactListParams{
		CustomerID: customerID,
		Page:       page,
		PageSize:   pageSize,
	}
	result, err := GetContactsByCustomerID(h.DBReplica, params)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "İletişim kayıtları alınamadı", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// İletişim kaydı ekleme (POST /api/contacts)
func (h *Handler) CreateContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		common.WriteError(w, http.StatusBadRequest, "Geçersiz istek gövdesi", err)
		return
	}
	if contact.CustomerID <= 0 {
		common.WriteError(w, http.StatusBadRequest, "Geçersiz müşteri ID", nil)
		return
	}
	if len(contact.Content) < 2 {
		common.WriteError(w, http.StatusBadRequest, "İletişim içeriği çok kısa", nil)
		return
	}
	if err := CreateContact(h.DBPrimary, &contact); err != nil {
		common.WriteError(w, http.StatusBadRequest, "İletişim kaydı eklenemedi", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}
