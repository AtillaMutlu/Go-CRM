package common

import (
	"net"
	"net/http"
	"strings"
	"time"

	redis_rate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var limiter *redis_rate.Limiter

// Rate limiter başlatıcı (singleton)
func InitRateLimiter(rdb *redis.Client) {
	limiter = redis_rate.NewLimiter(rdb)
}

// IP bazlı rate limiting middleware (örn. 5 istek/dk)
func RateLimitMiddleware(next http.Handler, limit int, period time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		rate := redis_rate.Limit{
			Rate:   limit,
			Period: period,
		}
		res, err := limiter.Allow(r.Context(), ip, rate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Rate limiter hatası"))
			return
		}
		if res.Allowed == 0 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Çok fazla istek, lütfen daha sonra tekrar deneyin."))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Gerçek istemci IP'sini bulur
func clientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
