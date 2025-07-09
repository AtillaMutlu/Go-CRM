package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	APIBaseURL = "http://localhost:8085"
	GatewayURL = "http://localhost:8080"
)

// Test kullanıcısı
var testUser = struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}{
	Email:    "demo@example.com",
	Password: "demo123",
}

// JWT token global değişkeni
var authToken string

// HTTP Client with timeout
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// Helper: HTTP request yapar
func makeRequest(method, url string, body interface{}, token string) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}
	
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	return client.Do(req)
}

// Test 1: Service Health Check
func TestServiceHealthChecks(t *testing.T) {
	t.Run("GatewayHealth", func(t *testing.T) {
		resp, err := client.Get(GatewayURL + "/healthz")
		if err != nil {
			t.Fatalf("Gateway health check başarısız: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Gateway health check status: %v", resp.StatusCode)
		}
		
		t.Log("✅ Gateway servis ayakta")
	})
}

// Test 2: Authentication Flow
func TestAuthenticationFlow(t *testing.T) {
	t.Run("ValidLogin", func(t *testing.T) {
		resp, err := makeRequest("POST", APIBaseURL+"/api/login", testUser, "")
		if err != nil {
			t.Fatalf("Login request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Login başarısız, status: %v", resp.StatusCode)
		}
		
		var loginResponse struct {
			Token string `json:"token"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
			t.Fatalf("Login response decode hatası: %v", err)
		}
		
		if loginResponse.Token == "" {
			t.Fatal("JWT token boş")
		}
		
		// Global token'ı set et
		authToken = loginResponse.Token
		t.Logf("✅ Login başarılı, token alındı: %.20s...", authToken)
	})
	
	t.Run("InvalidLogin", func(t *testing.T) {
		invalidUser := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    "invalid@example.com",
			Password: "wrongpassword",
		}
		
		resp, err := makeRequest("POST", APIBaseURL+"/api/login", invalidUser, "")
		if err != nil {
			t.Fatalf("Invalid login request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			t.Fatal("Invalid login başarılı olmamalıydı!")
		}
		
		t.Logf("✅ Invalid login doğru şekilde reddedildi, status: %v", resp.StatusCode)
	})
}

// Test 3: Customer Management E2E Flow
func TestCustomerManagementFlow(t *testing.T) {
	if authToken == "" {
		t.Skip("Auth token yok, önce authentication testini çalıştır")
	}
	
	var customerID int
	
	// Customer Creation
	t.Run("CreateCustomer", func(t *testing.T) {
		newCustomer := map[string]string{
			"name":  "E2E Test Müşteri",
			"email": "e2e@example.com",
			"phone": "+905551234567",
		}
		
		resp, err := makeRequest("POST", APIBaseURL+"/api/customers", newCustomer, authToken)
		if err != nil {
			t.Fatalf("Customer creation request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Customer oluşturulamadı, status: %v", resp.StatusCode)
		}
		
		var createdCustomer struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&createdCustomer); err != nil {
			t.Fatalf("Customer response decode hatası: %v", err)
		}
		
		if createdCustomer.ID == 0 {
			t.Fatal("Customer ID alınamadı")
		}
		
		customerID = createdCustomer.ID
		t.Logf("✅ Customer oluşturuldu: ID=%d, Name=%s", customerID, createdCustomer.Name)
	})
	
	// Customer List Retrieval
	t.Run("GetCustomers", func(t *testing.T) {
		resp, err := makeRequest("GET", APIBaseURL+"/api/customers", nil, authToken)
		if err != nil {
			t.Fatalf("Customer list request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Customer listesi alınamadı, status: %v", resp.StatusCode)
		}
		
		var customers []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&customers); err != nil {
			t.Fatalf("Customer list decode hatası: %v", err)
		}
		
		if len(customers) == 0 {
			t.Fatal("Customer listesi boş")
		}
		
		t.Logf("✅ %d customer bulundu", len(customers))
	})
	
	// Customer Update
	t.Run("UpdateCustomer", func(t *testing.T) {
		if customerID == 0 {
			t.Skip("Customer ID yok, create testi başarısız olmuş olabilir")
		}
		
		updatedCustomer := map[string]string{
			"name":  "Güncellenmiş E2E Müşteri",
			"email": "updated-e2e@example.com",
			"phone": "+905559876543",
		}
		
		url := fmt.Sprintf("%s/api/customers/%d", APIBaseURL, customerID)
		resp, err := makeRequest("PUT", url, updatedCustomer, authToken)
		if err != nil {
			t.Fatalf("Customer update request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Customer güncellenemedi, status: %v", resp.StatusCode)
		}
		
		t.Logf("✅ Customer güncellendi: ID=%d", customerID)
	})
	
	// Customer Deletion
	t.Run("DeleteCustomer", func(t *testing.T) {
		if customerID == 0 {
			t.Skip("Customer ID yok, create testi başarısız olmuş olabilir")
		}
		
		url := fmt.Sprintf("%s/api/customers/%d", APIBaseURL, customerID)
		resp, err := makeRequest("DELETE", url, nil, authToken)
		if err != nil {
			t.Fatalf("Customer delete request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			t.Fatalf("Customer silinemedi, status: %v", resp.StatusCode)
		}
		
		t.Logf("✅ Customer silindi: ID=%d", customerID)
	})
}

// Test 4: Contacts Retrieval
func TestContactsFlow(t *testing.T) {
	if authToken == "" {
		t.Skip("Auth token yok, önce authentication testini çalıştır")
	}
	
	t.Run("GetContacts", func(t *testing.T) {
		resp, err := makeRequest("GET", APIBaseURL+"/api/contacts", nil, authToken)
		if err != nil {
			t.Fatalf("Contacts request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Contacts alınamadı, status: %v", resp.StatusCode)
		}
		
		var contacts []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&contacts); err != nil {
			t.Fatalf("Contacts decode hatası: %v", err)
		}
		
		t.Logf("✅ %d contact bulundu", len(contacts))
	})
}

// Test 5: Authorization Tests
func TestAuthorizationFlow(t *testing.T) {
	t.Run("UnauthorizedAccess", func(t *testing.T) {
		// Token olmadan API'ye erişim dene
		resp, err := makeRequest("GET", APIBaseURL+"/api/customers", nil, "")
		if err != nil {
			t.Fatalf("Unauthorized request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			t.Fatal("Unauthorized request başarılı olmamalıydı!")
		}
		
		t.Logf("✅ Unauthorized access doğru şekilde reddedildi, status: %v", resp.StatusCode)
	})
	
	t.Run("InvalidToken", func(t *testing.T) {
		// Geçersiz token ile erişim dene
		resp, err := makeRequest("GET", APIBaseURL+"/api/customers", nil, "invalid-token")
		if err != nil {
			t.Fatalf("Invalid token request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			t.Fatal("Invalid token request başarılı olmamalıydı!")
		}
		
		t.Logf("✅ Invalid token doğru şekilde reddedildi, status: %v", resp.StatusCode)
	})
}

// Test 6: Error Handling
func TestErrorHandling(t *testing.T) {
	if authToken == "" {
		t.Skip("Auth token yok, önce authentication testini çalıştır")
	}
	
	t.Run("InvalidJSON", func(t *testing.T) {
		// Geçersiz JSON gönder
		invalidJSON := strings.NewReader(`{"invalid": json}`)
		req, _ := http.NewRequest("POST", APIBaseURL+"/api/customers", invalidJSON)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)
		
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Invalid JSON request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			t.Fatal("Invalid JSON request başarılı olmamalıydı!")
		}
		
		t.Logf("✅ Invalid JSON doğru şekilde reddedildi, status: %v", resp.StatusCode)
	})
	
	t.Run("NotFoundResource", func(t *testing.T) {
		// Var olmayan customer'ı güncellemeye çalış
		resp, err := makeRequest("PUT", APIBaseURL+"/api/customers/99999", map[string]string{
			"name": "Non-existent",
		}, authToken)
		if err != nil {
			t.Fatalf("Not found request hatası: %v", err)
		}
		defer resp.Body.Close()
		
		// 404 veya 500 bekleniyor (implementation'a göre)
		if resp.StatusCode == http.StatusOK {
			t.Log("⚠️  Non-existent resource update başarılı oldu (bu normal olabilir)")
		} else {
			t.Logf("✅ Non-existent resource update doğru şekilde handle edildi, status: %v", resp.StatusCode)
		}
	})
}

// Test 7: Performance & Load Test (Basic)
func TestBasicLoadTest(t *testing.T) {
	if authToken == "" {
		t.Skip("Auth token yok, önce authentication testini çalıştır")
	}
	
	t.Run("ConcurrentRequests", func(t *testing.T) {
		concurrency := 10
		done := make(chan bool, concurrency)
		errors := make(chan error, concurrency)
		
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				resp, err := makeRequest("GET", APIBaseURL+"/api/customers", nil, authToken)
				if err != nil {
					errors <- fmt.Errorf("Goroutine %d hatası: %v", id, err)
					return
				}
				defer resp.Body.Close()
				
				if resp.StatusCode != http.StatusOK {
					errors <- fmt.Errorf("Goroutine %d status hatası: %v", id, resp.StatusCode)
					return
				}
				
				done <- true
			}(i)
		}
		
		// Sonuçları topla
		successCount := 0
		errorCount := 0
		
		for i := 0; i < concurrency; i++ {
			select {
			case <-done:
				successCount++
			case err := <-errors:
				t.Logf("Error: %v", err)
				errorCount++
			case <-time.After(30 * time.Second):
				t.Fatal("Load test timeout")
			}
		}
		
		t.Logf("✅ Load test tamamlandı: %d başarılı, %d hatalı", successCount, errorCount)
		
		if successCount < concurrency/2 {
			t.Fatalf("Çok fazla hata: %d/%d", errorCount, concurrency)
		}
	})
} 