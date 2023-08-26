package orderbook

import list "github/wry-0313/exchange/linkedlist"

type OrderBook struct {
	orders map[string]*list.Node[*Order] // orderID -> *Order
	asks   *OrderSide
	bids   *OrderSide
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		orders: map[string]*list.Node[*Order]{},
		bids:   NewOrderSide(),
		asks:   NewOrderSide(),
	}
}
