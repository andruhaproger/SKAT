package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("Successfully connected to Redis!")
}

func SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return RedisClient.Set(Ctx, key, value, ttl).Err()
}

func Get(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}
