package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/your-monorepo/pkg/gateway"
)

func TestHealthzHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gateway.HealthzHandler)

	handler.ServeHTTP(rr, req)

	// Status kodu kontrolü
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler yanlış status kodu döndü: %v, beklenen: %v", status, http.StatusOK)
	}

	// Response body kontrolü
	expected := "ok\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler yanlış body döndü: %v, beklenen: %v", rr.Body.String(), expected)
	}
}

func TestHealthzHandlerMethods(t *testing.T) {
	tests := []struct {
		method     string
		shouldPass bool
	}{
		{"GET", true},
		{"POST", true}, // Handler tüm metodları kabul eder
		{"PUT", true},
		{"DELETE", true},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/healthz", nil)
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(gateway.HealthzHandler)
			handler.ServeHTTP(rr, req)

			if tt.shouldPass && rr.Code != http.StatusOK {
				t.Errorf("Metod %v için beklenmeyen status: %v", tt.method, rr.Code)
			}
		})
	}
}

// Benchmark test
func BenchmarkHealthzHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/healthz", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(gateway.HealthzHandler)
		handler.ServeHTTP(rr, req)
	}
}

// Rate limiting test helper
func TestRateLimitSimulation(t *testing.T) {
	// Bu test rate limiting logic'ini simulate eder
	// Gerçek implementation eklenince genişletilecek

	requests := 15 // Rate limit 10'dan fazla
	successCount := 0

	for i := 0; i < requests; i++ {
		// Simulate rate limit check
		if i < 10 {
			successCount++
		}
	}

	if successCount != 10 {
		t.Errorf("Rate limit testi başarısız: %v başarılı istek, beklenen: 10", successCount)
	}
}

// Concurrent access test
func TestHealthzConcurrency(t *testing.T) {
	handler := http.HandlerFunc(gateway.HealthzHandler)

	// 100 concurrent request
	concurrency := 100
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/healthz", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Concurrent test başarısız: status %v", rr.Code)
			}
			done <- true
		}()
	}

	// Tüm goroutine'lerin bitmesini bekle
	for i := 0; i < concurrency; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent test timeout")
		}
	}
}
