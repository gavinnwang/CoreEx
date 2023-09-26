package redis

import (
	"context"
	"fmt"
	"github/wry-0313/exchange/internal/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedis(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
	})

	// check redis connetcion
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Redis connected")
	}
	
	return rdb
}