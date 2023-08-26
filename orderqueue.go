package orderbook

import (
	"fmt"
	list "github/wry-0313/exchange/linkedlist"
	"strings"

	"github.com/shopspring/decimal"
)

// OrderQueue stores a queue of orders in a doubly linked list at a certain price level
type OrderQueue struct {
	volume decimal.Decimal
	price  decimal.Decimal
	orders *list.List[*Order]
}

// NewOrderQueue initializes a order queue of type orderbook.Order at a given price level. Defaults to zero total volume
func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{
		price:  price,
		volume: decimal.Zero,
		orders: list.New[*Order](),
	}
}

func (oq *OrderQueue) Len() int {
	return oq.orders.Len()
}

func (oq *OrderQueue) Price() decimal.Decimal {
	return oq.price
}

func (oq *OrderQueue) Volume() decimal.Decimal {
	return oq.volume
}

// Head returns a pointer to the Order at the front of the queue
func (oq *OrderQueue) Head() *list.Node[*Order] {
	return oq.orders.Front()
}

// Tail returns a poiner to the Order at the back of the queue
func (oq *OrderQueue) Tail() *list.Node[*Order] {
	return oq.orders.Back()
}

func (oq *OrderQueue) Append(o *Order) *list.Node[*Order] {
	oq.volume = oq.volume.Add(o.Volume())
	return oq.orders.PushBack(o)
}

func (oq *OrderQueue) Remove(n *list.Node[*Order]) *Order {
	oq.volume = oq.volume.Sub(n.Value.Volume())
	return oq.orders.Remove(n)
}

func (oq *OrderQueue) String() string {
	sb := strings.Builder{}
	iter := oq.orders.Front()
	sb.WriteString(fmt.Sprintf("\nqueue length: %d, price: %s, volume: %s, orders:", oq.Len(), oq.Price(), oq.Volume()))
	for iter != nil {
		order := iter.Value
		str := fmt.Sprintf("\n\tid: %s, volume: %s, time: %s", order.OrderID(), order.Volume(), order.Price())
		sb.WriteString(str)
		iter = iter.Next()
	}
	return sb.String()
}
