package models

import "time"

type Stock struct {
	ID          int       `json:"id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
}

type StockPriceHistory struct {
	StockID    int     `json:"stock_id"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	RecordedAt int64   `json:"recorded_at"`
}

type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	StockID     int       `json:"stock_id"`
	OrderType   string    `json:"order_type"`
	OrderStatus string    `json:"order_status"`
	Quantity    float64   `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Transaction struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	StockID         int     `json:"stock_id"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	TransactionType string  `json:"transaction_type"`
}
