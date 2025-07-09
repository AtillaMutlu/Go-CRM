package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Go-CRM/pkg/customer"
)

func TestCreateCustomerHandler_InvalidEmail(t *testing.T) {
	h := &customer.Handler{DBPrimary: nil, DBReplica: nil}
	payload := map[string]interface{}{
		"name":  "Ali",
		"email": "yanlisemail",
		"phone": "5551234567",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateCustomerHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Beklenen 400, gelen %d", w.Code)
	}
}
