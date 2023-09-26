package exchange

type PlaceOrderInput struct {
	UserID    string  `json:"user_id" validate:"omitempty"`
	OrderType string  `json:"order_type" validate:"required,oneof=market limit"`
	OrderSide string  `json:"order_side" validate:"required,oneof=buy sell"`
	Price     float64 `json:"price" validate:"required_if=OrderType limit"`
	Volume    float64 `json:"volume" validate:"required"`
	Symbol    string  `json:"symbol" validate:"required"`
}
