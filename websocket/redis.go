package ws

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

func NewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb
}
