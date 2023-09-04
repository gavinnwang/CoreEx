package exchange

type PlaceOrderInput struct {
	ClientID  string  `json:"client_id"`
	OrderType string  `json:"order_type"`
	OrderSide string  `json:"order_side"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Symbol    string  `json:"symbol"`
}

type StreamPriceInput struct {
	Symbol string `json:"symbol"`
}