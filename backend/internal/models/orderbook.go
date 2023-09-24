package models

import "github.com/shopspring/decimal"

type Stock struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
}

type StockPriceHistory struct {
	Symbol     string  `json:"symbol"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	RecordedAt int64   `json:"recorded_at"`
}

type Holding struct {
	UserID       string  `json:"user_id"`
	Symbol       string  `json:"symbol"`
	VolumeChange decimal.Decimal `json:"volume_change"`
}
