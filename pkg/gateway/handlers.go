package gateway

import (
	"fmt"
	"net/http"
)

// /healthz endpoint handler'ı
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}
