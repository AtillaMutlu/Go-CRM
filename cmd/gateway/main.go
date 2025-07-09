package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"sync"

	"Go-CRM/pkg/gateway"
)

// Rate limit için basit sayaç
var rateLimiters = make(map[string]*rateLimiter)
var rlMu sync.Mutex

// IP başına rate limiter
type rateLimiter struct {
	count     int
	timestamp time.Time
}

// Allowlist IP'ler
var allowlist = map[string]bool{
	"127.0.0.1": true,
	"::1":       true,
}

func main() {
	/*
		var err error
		jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
			RefreshInterval: time.Minute * 5,
			RefreshErrorHandler: func(err error) {
				log.Printf("JWKS yenileme hatası: %v", err)
			},
			RefreshTimeout:    10 * time.Second,
			RefreshUnknownKID: true,
		})
		if err != nil {
			log.Fatalf("JWKS alınamadı: %v", err)
		}
	*/

	http.HandleFunc("/healthz", gateway.HealthzHandler)

	http.HandleFunc("/api", ipAllowlist(rateLimit(apiHandler)))

	fmt.Println("Gateway servis 8080 portunda başlatıldı...")
	http.ListenAndServe(":8080", nil)
}

// Rate limit middleware
func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if i := strings.LastIndex(ip, ":"); i != -1 {
			ip = ip[:i]
		}
		rlMu.Lock()
		rl, ok := rateLimiters[ip]
		if !ok || time.Since(rl.timestamp) > time.Minute {
			rl = &rateLimiter{count: 1, timestamp: time.Now()}
			rateLimiters[ip] = rl
		} else {
			rl.count++
		}
		count := rl.count
		rlMu.Unlock()
		if count > 10 {
			http.Error(w, "Çok fazla istek! (rate limit)", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

// IP allowlist middleware
func ipAllowlist(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if i := strings.LastIndex(ip, ":"); i != -1 {
			ip = ip[:i]
		}
		if !allowlist[ip] {
			http.Error(w, "IP erişimine izin verilmiyor", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// Korunan örnek endpoint
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "JWT dogrulamasi gecici olarak devre disi! Korumali /api endpointindesin.")
}
