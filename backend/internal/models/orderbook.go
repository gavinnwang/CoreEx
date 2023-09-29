package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Stock struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
}

type StockPriceHistory struct {
	PriceData
	BidVolume  float64 `json:"bid_volume"`
	AskVolume  float64 `json:"ask_volume"`
	RecordedAt time.Time   `json:"recorded_at"`
}

type PriceData struct {
	Open  decimal.Decimal `json:"open"`
	Close decimal.Decimal `json:"close"`
	High  decimal.Decimal `json:"high"`
	Low   decimal.Decimal `json:"low"`
}

type HoldingChange struct {
	UserID       string          `json:"user_id"`
	Symbol       string          `json:"symbol"`
	VolumeChange decimal.Decimal `json:"volume_change"`
}

type Holding struct {
	Symbol string          `json:"symbol"`
	Volume decimal.Decimal `json:"volume"`
}

type Order struct {
	Symbol         string    `json:"symbol"`
	OrderID        string    `json:"order_id"`
	OrderSide      string    `json:"order_side"`
	OrderStatus    string    `json:"order_status"`
	OrderType      string    `json:"order_type"`
	FilledAt       float64   `json:"filled_at"`
	FilledAtTime   time.Time `json:"filled_at_time"`
	TotalProcessed float64   `json:"total_processed"`
	Volume         float64   `json:"volume"`
	InitialVolume  float64   `json:"initial_volume"`
	Price          float64   `json:"price"`
}
