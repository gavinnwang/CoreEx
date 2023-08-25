package orderbook

import (
	list "github/wry-0313/exchange/linkedlist"
	"github/wry-0313/exchange/treemap"

	"github.com/shopspring/decimal"
)

type OrderSide struct {
	// a sorted treemap with price being keys and order queues as values
	priceTree  *treemap.TreeMap[decimal.Decimal, *OrderQueue]
	priceTable map[string]*OrderQueue

	volume    decimal.Decimal
	depth     int
	numOrders int
}

func keyComparator(a, b decimal.Decimal) bool {
	return a.Cmp(b) == -1
}

func NewOrderSide() *OrderSide {
	return &OrderSide{
		priceTree:  treemap.NewWith[decimal.Decimal, *OrderQueue](keyComparator),
		priceTable: map[string]*OrderQueue{},
		volume:     decimal.Zero,
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
	// if priceQueue at price level doesn't exit
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
func (os *OrderSide) MaxPriceQueue() *OrderQueue {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMax(); found {
			return value
		}
	}
	return nil
}

// MinPriceQueue returns minimal level of price
func (os *OrderSide) MinPriceQueue() *OrderQueue {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMin(); found {
			return value
		}
	}
	return nil
}
