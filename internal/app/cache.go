package app

import (
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	*marshaler.Marshaler
}

func NewCacheService() *CacheService {
	cacheStore := redis_store.NewRedis(redis.NewClient(redisOptions()))
	cacheManager := cache.New[any](cacheStore)
	m := marshaler.New(cacheManager)
	return &CacheService{
		m,
	}
}
