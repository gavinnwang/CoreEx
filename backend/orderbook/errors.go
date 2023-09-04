package orderbook

import "errors"

var (
	ErrInvalidVolume   = errors.New("orderbook: invalid order volume")
	ErrInvalidClientID = errors.New("orderbook: invalid client ID")
	ErrInvalidPrice    = errors.New("orderbook: invalid order price")
	ErrInvalidSide     = errors.New("orderbook: invalid order side")
	ErrOrderExists     = errors.New("orderbook: order already exists")
	ErrOrderNotExists  = errors.New("orderbook: order does not exist")
)
