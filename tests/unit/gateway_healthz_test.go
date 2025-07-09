package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/your-monorepo/pkg/gateway"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	gateway.HealthzHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Beklenen status 200, gelen: %d", resp.StatusCode)
	}
}
