package main

import (
	"fmt"
	"time"
)

type Order struct {
	Size float64
	Bid bool
	Limit *Limit
	Timestamp int64
}

func NewOrder(bid bool, size float64) *Order {
	return &Order {
		Size: size,
		Bid: bid,
		Timestamp: time.Now().UnixNano(),	
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

// Group of orders at a certain price limit 
type Limit struct {
	Price float64
	Orders []*Order
	TotalVolume float64
}

func NewLimit(price float64) *Limit {
	return &Limit { 
		Price: price,
		Orders: []*Order{}, // {} dentoes an empty slice literal 
	}
}

// (l *Limit) is a receiver, it's a method on the Limit struct
func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit 
}