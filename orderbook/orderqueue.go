package orderbook

import (
	"fmt"
	list "github/wry-0313/exchange/linkedlist"
	"strings"
	"sync"

	// "sync"

	"github.com/shopspring/decimal"
)

// OrderQueue stores a queue of orders in a doubly linked list at a certain price level
type OrderQueue struct {
	volume   decimal.Decimal    // volume can be changed so we need a mutex
	volumeMu sync.RWMutex       // protect volume
	price    decimal.Decimal    // price level cannot be changed once initialized
	orders   *list.List[*Order] // limit orders
	ordersMu sync.RWMutex       // protect orders
}

// NewOrderQueue initializes a order queue of type orderbook.Order at a given price level. Defaults to zero total volume
func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{
		price:    price,
		volume:   decimal.Zero,
		orders:   list.New[*Order](),
	}
}

func (oq *OrderQueue) Len() int {
	oq.ordersMu.RLock()
	defer oq.ordersMu.RUnlock()
	return oq.orders.Len()
}

func (oq *OrderQueue) Price() decimal.Decimal {
	return oq.price
}

// Head returns a pointer to the Order at the front of the queue
func (oq *OrderQueue) Head() *list.Node[*Order] {
	oq.ordersMu.RLock()
	defer oq.ordersMu.RUnlock()
	return oq.orders.Front()
}

// Tail returns a poiner to the Order at the back of the queue
// func (oq *OrderQueue) Tail() *list.Node[*Order] {
// 	oq.ordersMu.RLock()
// 	defer oq.ordersMu.RUnlock()
// 	return oq.orders.Back()
// }

func (oq *OrderQueue) Volume() decimal.Decimal {
	oq.volumeMu.RLock()
	defer oq.volumeMu.RUnlock()
	return oq.volume
}

func (oq *OrderQueue) SetVolume(volume decimal.Decimal) decimal.Decimal {
	oq.volumeMu.Lock()
	oq.volume = volume
	oq.volumeMu.Unlock()
	return volume
}

func (oq *OrderQueue) Append(o *Order) *list.Node[*Order] {
	oq.volumeMu.Lock()
	oq.volume = oq.volume.Add(o.Volume())
	oq.volumeMu.Unlock()
	oq.ordersMu.Lock()
	defer oq.ordersMu.Unlock()
	return oq.orders.PushBack(o)
}

func (oq *OrderQueue) Remove(n *list.Node[*Order]) *Order {
	oq.volumeMu.Lock()
	oq.volume = oq.volume.Sub(n.Value.Volume())
	oq.volumeMu.Unlock()
	oq.ordersMu.Lock()
	defer oq.ordersMu.Unlock()
	return oq.orders.Remove(n)
}

func (oq *OrderQueue) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("queue: length: %d, price: %s, volume: %s\n", oq.Len(), oq.Price(), oq.Volume()))

	// iter := oq.orders.Front()
	// for iter != nil {
	// 	order := iter.Value
	// 	str := fmt.Sprintf("\n\tid: %s, volume: %s, price: %s", order.OrderID(), order.Volume(), order.Price())
	// 	sb.WriteString(str)
	// 	iter = iter.Next()
	// }
	return sb.String()
}
