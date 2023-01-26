package app

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func redisOptions() *redis.Options {
	opt := &redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	}

	username := os.Getenv("REDIS_USERNAME")
	if username != "" {
		opt.Username = username
	}

	passwd := os.Getenv("REDIS_PASSWORD")
	if passwd != "" {
		opt.Password = passwd
	}

	return opt
}
