package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
)

func Connect() {
	Redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
