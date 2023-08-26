package orderbook

import "errors"

var (
	ErrInvalidVolume      = errors.New("orderbook: invalid order volume")
	ErrInvalidPrice         = errors.New("orderbook: invalid order price")
	ErrOrderExists          = errors.New("orderbook: order already exists")
	ErrOrderNotExists       = errors.New("orderbook: order does not exist")
)