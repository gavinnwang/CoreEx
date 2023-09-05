package ws

import (
	"context"
	"fmt"
	"github/wry-0313/exchange/config"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

func NewRedis(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
	})
	return rdb
}
