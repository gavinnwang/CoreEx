package models

type Stock struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
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
	StockID     int     `json:"stock_id"`
	OrderType   string  `json:"order_type"`
	OrderStatus string  `json:"order_status"`
	Volume      float64 `json:"volume"`
	Price       float64 `json:"price"`
	OrderSide   string  `json:"order_side"`
}

type Transaction struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	StockID         int     `json:"stock_id"`
	Volume          float64 `json:"volume"`
	Price           float64 `json:"price"`
	TransactionType string  `json:"transaction_type"`
}
