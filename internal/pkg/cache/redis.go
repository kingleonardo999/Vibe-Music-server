package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
	"vibe-music-server/internal/config"
)

var (
	cache             *redis.Client
	once              sync.Once
	defaultExpiration time.Duration
)

func Init() {
	once.Do(func() {
		Redis := config.Get().Redis
		cache = redis.NewClient(&redis.Options{
			Addr:     Redis.Host + ":" + Redis.Port,
			Password: Redis.Password,
			DB:       Redis.Database,
		})
		defaultExpiration = time.Duration(Redis.TimeToLive) * time.Second
	})
}

func Cache() *redis.Client {
	return cache
}

// SetWithExp 设置缓存, 带过期时间
func SetWithExp(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(cache.Context(), 3*time.Second)
	defer cancel()
	return cache.Set(ctx, key, value, expiration).Err()
}

func Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(cache.Context(), 3*time.Second)
	defer cancel()
	return cache.Set(ctx, key, value, defaultExpiration).Err()
}

func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(cache.Context(), 3*time.Second)
	defer cancel()
	return cache.Get(ctx, key).Result()
}

func Del(key string) error {
	ctx, cancel := context.WithTimeout(cache.Context(), 3*time.Second)
	defer cancel()
	return cache.Del(ctx, key).Err()
}
