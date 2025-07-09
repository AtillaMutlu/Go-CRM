package common

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// Redis bağlantısını başlatır (singleton)
func InitRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Auth yoksa boş bırak
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return redisClient.Ping(ctx).Err()
}

// Redis'ten string değer okur
func RedisGet(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

// Redis'e string değer yazar (TTL ile)
func RedisSet(ctx context.Context, key, value string, ttl time.Duration) error {
	return redisClient.Set(ctx, key, value, ttl).Err()
}

// Redis client getter
func GetRedisClient() *redis.Client {
	return redisClient
}
