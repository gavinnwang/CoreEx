package exchange

type PlaceOrderInput struct {
	UserID    string  `json:"user_id" validate:"omitempty"`
	OrderType string  `json:"order_type" validate:"required,oneof=market limit"`
	OrderSide string  `json:"order_side" validate:"required,oneof=buy sell"`
	Price     float64 `json:"price" validate:"required_if=OrderType limit"`
	Volume    float64 `json:"volume" validate:"required"`
	Symbol    string  `json:"symbol" validate:"required"`
}

type StreamPriceParams struct {
	Symbol string `json:"symbol" validate:"required"`
}

type SymbolInfoResponse struct {
	Symbol    string  `json:"symbol"`
	AskVolume float64 `json:"ask_volume"`
	BidVolume float64 `json:"bid_volume"`
	BestBid   float64 `json:"best_bid"`
	BestAsk   float64 `json:"best_ask"`
	Price     float64 `json:"price"`
}