package ws

import (
	"github/wry-0313/exchange/internal/exchange"

	"github.com/go-redis/redis"
)

type WebSocket struct {
	exchangeService exchange.Service
	rdb             *redis.Client
}

func NewWebSocket(exchangeService exchange.Service, rdb *redis.Client) *WebSocket {
	return &WebSocket{
		exchangeService: exchangeService,
		rdb:             rdb,
	}
}
