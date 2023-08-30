package orderbook

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type OrderLogs struct {
	entries []LogEntry
}

type LogEntry struct {
	timestamp time.Time
	message   string
}

func NewLimitOrderLogs(volume, price decimal.Decimal) *OrderLogs {
	return new(OrderLogs).initLimit(volume, price)
}

func NewMarketOrderLogs(volume decimal.Decimal) *OrderLogs {
	return new(OrderLogs).initMarket(volume)
}

func (ol *OrderLogs) initLimit(volume, price decimal.Decimal) *OrderLogs {
	ol.Log(fmt.Sprintf("Initialized limit order with volume: %v, price: %v", volume, price))
	return ol
}
func (ol *OrderLogs) initMarket(volume decimal.Decimal) *OrderLogs {
	ol.Log(fmt.Sprintf("Initialized market order with volume %v", volume))
	return ol
}

func (ol *OrderLogs) Log(logMsg string) {
	ol.entries = append(ol.entries, LogEntry{
		timestamp: time.Now(),
		message:   logMsg,
	})
}
