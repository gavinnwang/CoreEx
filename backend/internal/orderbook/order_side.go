package orderbook

import (
	"fmt"
	list "github/wry-0313/exchange/pkg/dsa/linkedlist"

	// "github/wry-0313/exchange/internaltreemap"
	"strings"
	"sync"

	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"

	"github.com/shopspring/decimal"
)

type OrderSide struct {
	priceTree  *rbtx.RedBlackTreeExtended // price -> *OrderQueue, sorted by price
	priceTable map[string]*OrderQueue     // price -> *OrderQueue for quick lookup

	volume   decimal.Decimal // total volume of all orders
	volumeMu sync.RWMutex    // protect volume

	depth   int          // number of active price levels
	depthMu sync.RWMutex // protect depth

	numOrders   int          // number of orders
	numOrdersMu sync.RWMutex // protect numOrders
}

func keyComparator(a, b decimal.Decimal) bool {
	return a.Cmp(b) == -1
}

func NewOrderSide() *OrderSide {
	return &OrderSide{
		priceTree: &rbtx.RedBlackTreeExtended{
			Tree: rbt.NewWith(rbtComparator),
		},
		priceTable: map[string]*OrderQueue{},
		volume:     decimal.Zero,
		depth:      0,
		numOrders:  0,
	}
}

func rbtComparator(a, b interface{}) int {
	return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
}

func (os *OrderSide) Len() int {
	os.numOrdersMu.RLock()
	defer os.numOrdersMu.RUnlock()
	return os.numOrders
}

func (os *OrderSide) Depth() int {
	os.depthMu.RLock()
	defer os.depthMu.RUnlock()
	return os.depth
}

func (os *OrderSide) Append(o *Order) *list.Node[*Order] {
	price := o.Price()
	priceStr := price.String()

	// os.priceTreeMu.Lock()
	// defer os.priceTreeMu.Unlock()

	// os.priceTableMu.RLock()
	priceQueue, ok := os.priceTable[priceStr]
	// os.priceTableMu.RUnlock()
	// if priceQueue at price level doesn't exit, create a new order queue at that order level
	if !ok {
		priceQueue = NewOrderQueue(o.Price())
		// os.priceTableMu.Lock()
		os.priceTable[priceStr] = priceQueue
		// os.priceTableMu.Unlock()

		os.priceTree.Put(price, priceQueue)

		os.depthMu.Lock()
		os.depth++
		os.depthMu.Unlock()
	}
	os.numOrdersMu.Lock()
	os.numOrders++
	os.numOrdersMu.Unlock()
	// os.volume = os.volume.Add(o.Volume())
	// os.AddVolumeBy(o.Volume())
	return priceQueue.Append(o)
}

// Time Complexity: O(1) if don't remove price queue O(N) otherwise
func (os *OrderSide) Remove(n *list.Node[*Order]) *Order {
	price := n.Value.Price()
	priceStr := price.String()

	// os.priceTableMu.RLock()
	priceQueue, found := os.priceTable[priceStr]
	if !found {
		// Log(fmt.Sprintf("already removed: %s\n", priceQueue))
		return n.Value
	}
	// os.priceTableMu.RUnlock()

	// os.priceTreeMu.Lock()
	// defer os.priceTreeMu.Unlock()

	o := priceQueue.Remove(n)

	if priceQueue.Len() == 0 {
		// Log(fmt.Sprintf("Remove price queue at price level %s", priceStr))
		// os.priceTableMu.Lock()
		delete(os.priceTable, priceStr)
		// os.priceTableMu.Unlock()

		os.priceTree.Remove(price)
		// if !removed {
		// 	Log(fmt.Sprintf("Error: price level not removed from tree at price level %s", priceStr))
		// 	panic("price level not removed from tree")
		// }

		// Log(fmt.Sprintf("price level removed from tree at price level %s", priceStr))

		os.depthMu.Lock()
		os.depth--
		os.depthMu.Unlock()
	}

	os.numOrdersMu.Lock()
	os.numOrders--
	os.numOrdersMu.Unlock()

	// os.SubVolumeBy(o.Volume())

	return o
}

// MaxPriceQueue returns maximal level of price
func (os *OrderSide) MaxPriceQueue() (*OrderQueue, bool) {
	if os.Depth() > 0 {
		// os.priceTreeMu.Lock()

		if oq, found := os.priceTree.GetMax(); found {
			// Log(fmt.Sprintf("maxqueu len: %d\n", oq.(*OrderQueue).Len()))
			if oq.(*OrderQueue).Len() == 0 {
				// Log(fmt.Sprintf("Error: MaxPriceQueue: price queue is empty: %s\n", oq))
				// Log(fmt.Sprintf("max price error: os: %s\n", os))
				// max, _ := os.priceTree.GetMax()
				// Log(max.(*OrderQueue).String())
				panic("MaxPriceQueue: price queue is empty")
			}
			// os.priceTreeMu.Unlock()
			return oq.(*OrderQueue), true
		}
	}
	// os.priceTreeMu.Unlock()
	return nil, false
}

// MinPriceQueue returns minimal level of price
func (os *OrderSide) MinPriceQueue() (*OrderQueue, bool) {
	if os.Depth() > 0 {
		// os.priceTreeMu.RLock()
		// defer os.priceTreeMu.RUnlock()
		if oq, found := os.priceTree.GetMin(); found {
			// Log(fmt.Sprintf("MinPriceQueue len: %d\n", oq.(*OrderQueue).Len()))
			if oq.(*OrderQueue).Len() == 0 {
				// Log(fmt.Sprintf("Error: MinPriceQueue: price queue is empty: %s\n", oq))
				// Log(fmt.Sprintf("min price error: os: %s\n", os))
				// min, _ := os.priceTree.GetMin()
				// Log(min.(*OrderQueue).String())
				panic("Min PriceQueue: price queue is empty")
			}
			return oq.(*OrderQueue), true
		}
	}
	return nil, false
}

// func (os *OrderSide) SetVolume(volume decimal.Decimal) {
// 	os.volumeMu.Lock()
// 	defer os.volumeMu.Unlock()
// 	os.volume = volume
// }

func (os *OrderSide) AddVolumeBy(volume decimal.Decimal) {
	os.volumeMu.Lock()
	defer os.volumeMu.Unlock()
	os.volume = os.volume.Add(volume)
}

func (os *OrderSide) ResetVolume() {
	os.volumeMu.Lock()
	defer os.volumeMu.Unlock()
	os.volume = decimal.Zero
}

// func (os *OrderSide) SubVolumeBy(volume decimal.Decimal) {
// 	os.volumeMu.Lock()
// 	defer os.volumeMu.Unlock()
// 	os.volume = os.volume.Sub(volume)
// }



func (os *OrderSide) Volume() decimal.Decimal {
	os.volumeMu.RLock()
	defer os.volumeMu.RUnlock()
	return os.volume
}

// LessThan returns nearest OrderQueue with price less than given
// deprecate this and replace with iter()
// TODO
func (os *OrderSide) LessThan(price decimal.Decimal) *OrderQueue {
	tree := os.priceTree.Tree
	node := tree.Root

	var floor *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) > 0 {
			floor = node
			node = node.Right
		} else {
			node = node.Left
		}
	}

	if floor != nil {
		return floor.Value.(*OrderQueue)
	}

	return nil
}

func (os *OrderSide) String() string {
	sb := strings.Builder{}

	level, _ := os.MaxPriceQueue()
	for level != nil {
		sb.WriteString(fmt.Sprintf("\n%s -> %s", level.Price(), level.Volume()))
		level = os.LessThan(level.Price())
	}

	return sb.String()
}
