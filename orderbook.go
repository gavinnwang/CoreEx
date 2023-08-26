package orderbook

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	orders       map[uuid.UUID]*Order // orderID -> *Order
	currentPrice decimal.Decimal
	asks         *OrderSide
	bids         *OrderSide
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		orders: map[uuid.UUID]*Order{},
		bids:   NewOrderSide(),
		asks:   NewOrderSide(),
	}
}

func (ob *OrderBook) FillMarketOrder(side Side, clientID uuid.UUID, volume decimal.Decimal) {
	// create the order and add it to all orders
	o := NewOrder(side, clientID, Market, decimal.Zero, volume, true)
	ob.orders[o.OrderID()] = o
	
	// fill the order by taking the best limit orders of the opposite side
	// 
	if side == Buy {
		queue := ob.asks.MinPriceQueue();
		
	}
}
