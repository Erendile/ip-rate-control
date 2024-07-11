package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"ip-rate-control/pkg/config"
	"log"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

func InitializeRedis(cfg config.RedisConfig) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return client
}
