package common

import (
	"encoding/json"
	"log"
	"net/http"
)

// Ortak hata yazıcı. Kullanıcıya sade mesaj, log'a detaylı hata basar.
func WriteError(w http.ResponseWriter, status int, userMessage string, devError error) {
	log.Printf("[ERROR] %s | %v", userMessage, devError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": userMessage,
	})
}
