package models

import "github.com/shopspring/decimal"

type Stock struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
}

type StockPriceHistory struct {
	PriceData
	Volume     decimal.Decimal `json:"volume"`
	RecordedAt int64           `json:"recorded_at"`
}

type PriceData struct {
	Open  decimal.Decimal `json:"open"`
	Close decimal.Decimal `json:"close"`
	High  decimal.Decimal `json:"high"`
	Low   decimal.Decimal `json:"low"`
}

type Holding struct {
	UserID       string          `json:"user_id"`
	Symbol       string          `json:"symbol"`
	VolumeChange decimal.Decimal `json:"volume_change"`
}