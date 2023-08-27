package orderbook

import (
	list "github/wry-0313/exchange/linkedlist"
	"github/wry-0313/exchange/treemap"

	"github.com/shopspring/decimal"
)

type OrderSide struct {
	priceTree  *treemap.TreeMap[decimal.Decimal, *OrderQueue] // price -> *OrderQueue, sorted by price
	priceTable map[string]*OrderQueue                         // price -> *OrderQueue for quick lookup
	volume     decimal.Decimal                                // total volume of all orders
	depth      int                                            // number of price levels
	numOrders  int                                            // number of orders
}

func keyComparator(a, b decimal.Decimal) bool {
	return a.Cmp(b) == -1
}

func NewOrderSide() *OrderSide {
	return &OrderSide{
		priceTree:  treemap.NewWith[decimal.Decimal, *OrderQueue](keyComparator),
		priceTable: map[string]*OrderQueue{},
		volume:     decimal.Zero,
		depth:      0,
		numOrders:  0,
	}
}

func (os *OrderSide) Len() int {
	return os.numOrders
}

func (os *OrderSide) Depth() int {
	return os.depth
}

func (os *OrderSide) Volume() decimal.Decimal {
	return os.volume
}

func (os *OrderSide) Append(o *Order) {
	price := o.Price()
	priceStr := price.String()

	priceQueue, ok := os.priceTable[priceStr]
	// if priceQueue at price level doesn't exit, create a new order queue at that order level
	if !ok {
		priceQueue = NewOrderQueue(o.Price())
		os.priceTable[priceStr] = priceQueue
		os.priceTree.Add(price, priceQueue)
		os.depth++
	}
	os.numOrders++
	os.volume = os.volume.Add(o.Volume())
	priceQueue.Append(o)
}

func (os *OrderSide) Remove(n *list.Node[*Order]) *Order {
	price := n.Value.Price()
	priceStr := price.String()

	priceQueue := os.priceTable[priceStr]
	o := priceQueue.Remove(n)

	if priceQueue.Len() == 0 {
		delete(os.priceTable, priceStr)
		os.priceTree.Remove(price)
		os.depth--
	}

	os.numOrders--
	os.volume = os.volume.Sub(o.Volume())
	return o
}

// MaxPriceQueue returns maximal level of price
func (os *OrderSide) MaxPriceQueue() (*OrderQueue, bool) {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMax(); found {
			return value, true
		}
	}
	return nil, false
}

// MinPriceQueue returns minimal level of price
func (os *OrderSide) MinPriceQueue() (*OrderQueue, bool) {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMin(); found {
			return value, true
		}
	}
	return nil, false
}
