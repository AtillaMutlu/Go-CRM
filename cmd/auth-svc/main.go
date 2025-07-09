package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Auth servis çalışıyor. Keycloak proxy/adapter mantığı ile ilerleyecek.")
	})

	fmt.Println("Auth servis 8082 portunda başlatıldı...")
	http.ListenAndServe(":8082", nil)
}
