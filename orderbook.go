package orderbook

import (
	// "fmt"
	"fmt"
	list "github/wry-0313/exchange/linkedlist"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	activeOrders     map[uuid.UUID]*list.Node[*Order] // orderID -> *Order for quick acctions such as update or cancel
	ordersMu         sync.RWMutex
	marketPrice      decimal.Decimal
	marketPriceMu    sync.RWMutex
	marketBuyOrders  *list.List[*Order] // partially filled market buy orders
	marketSellOrders *list.List[*Order] // partially filled market sell orders
	asks             *OrderSide         // limit sell orders
	bids             *OrderSide         // limit buy orders
}

func NewOrderBook() *OrderBook {
	err := InitializeLogService("orderbook_log.txt")
	if err != nil {
		log.Fatalf("Could not initialize log service: %v", err)
	}
	return &OrderBook{
		activeOrders:     map[uuid.UUID]*list.Node[*Order]{},
		bids:             NewOrderSide(),
		asks:             NewOrderSide(),
		marketBuyOrders:  list.New[*Order](),
		marketSellOrders: list.New[*Order](),
		marketPrice:      decimal.Zero,
	}
}

func (ob *OrderBook) PlaceMarketOrder(side Side, clientID uuid.UUID, volume decimal.Decimal) (orderID uuid.UUID, err error) {
	if volume.Sign() <= 0 {
		return uuid.Nil, ErrInvalidVolume
	}
	// create the order and add it to all orders
	o := NewOrder(side, clientID, Market, decimal.Zero, volume, true)
	Log(fmt.Sprintf("Created market order: %v", o))
	// fill the order by taking the best limit orders of the opposite Side
	var (
		os   *OrderSide
		iter func() (*OrderQueue, bool)
	)
	if side == Buy {
		iter = ob.asks.MinPriceQueue
		os = ob.asks
	} else {
		iter = ob.bids.MaxPriceQueue
		os = ob.bids
	}

	if os.Len() == 0 {
		// no limit orders in the opposite side, add the order to the market order list
		ob.addMarketOrder(o)
		Log(fmt.Sprintf("No limit orders in the opposite side, add the order to the market order list: %s", o.shortOrderID()))
		return o.orderID, nil
	}

	var volumeLeft = o.Volume()

	for volumeLeft.Sign() > 0 && os.Len() > 0 { // while the order is not fully filled and the opposite side has more limit orders
		oq, _ := iter()
		volumeLeft = ob.matchAtPriceLevel(oq, o)
	}

	if volumeLeft.Sign() > 0 {
		// the order is not fully filled, add it to the market order list
		Log(fmt.Sprintf("The order is not fully filled, add it to the market order list: %s", o.shortOrderID()))
		ob.addMarketOrder(o)
	}
	return o.orderID, nil
}

func (ob *OrderBook) matchAtPriceLevel(oq *OrderQueue, o *Order) (volumeLeft decimal.Decimal) {
	volumeLeft = o.Volume()

	Log(fmt.Sprintf("Matching %s at price level %s", o.shortOrderID(), oq.Price()))

	for oq.Len() > 0 && volumeLeft.Sign() > 0 { // while the price level have more limit orders and the order is not fully filled

		bestOrderNode := oq.Head()
		bestOrder := oq.Head().Value

		Log(fmt.Sprintf("Matching %s with %s", o.shortOrderID(), bestOrder.shortOrderID()))

		if volumeLeft.LessThan(bestOrder.Volume()) { // the best order will be partially filled

			Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s", o.shortOrderID(), o.Volume(), o.Volume().Sub(volumeLeft), bestOrder.shortOrderID(), bestOrder.Volume(), bestOrder.Volume().Sub(volumeLeft)))

			matchedVolumeLeft := bestOrder.Volume().Sub(volumeLeft) // update order status. This change should reflect in order queue
			oq.volume = oq.volume.Sub(volumeLeft)

			if o.Side() == Buy {
				ob.bids.volume = ob.bids.volume.Sub(volumeLeft)
			} else {
				ob.asks.volume = ob.asks.volume.Sub(volumeLeft)
			}

			bestOrder.setStatusToPartiallyFilled(matchedVolumeLeft)
			o.setStatusToFilled()

			volumeLeft = decimal.Zero
			// Log(fmt.Sprintf("%s is partially filled and %s is completely filled", bestOrder.shortOrderID(), o.shortOrderID()))
		} else { // the best order will be completely filled
			volumeLeft = volumeLeft.Sub(bestOrder.Volume())

			Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s", o.shortOrderID(), o.Volume(), o.Volume().Sub(bestOrder.Volume()), bestOrder.shortOrderID(), bestOrder.Volume(), decimal.Zero))

			ob.fillAndRemoveLimitOrder(bestOrderNode)
			o.setStatusToPartiallyFilled(volumeLeft)
		}
	}
	return
}

func (ob *OrderBook) matchWithMarketOrders(marketOrders *list.List[*Order], order *Order) {
	for marketOrders.Len() > 0 {

		marketOrderNode := marketOrders.Front()
		marketOrder := marketOrderNode.Value
		marketOrderVolume := marketOrder.Volume()
		orderVolume := order.Volume()

		if orderVolume.LessThan(marketOrderVolume) { // the market order will be completely filled

			Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s", order.shortOrderID(), order.Volume(), decimal.Zero, marketOrder.shortOrderID(), marketOrder.Volume(), marketOrder.Volume().Sub(orderVolume)))

			marketOrder.setStatusToPartiallyFilled(marketOrderVolume.Sub(orderVolume))

			if order.Side() == Buy {
				ob.asks.volume = ob.asks.volume.Sub(orderVolume)
			} else {
				ob.bids.volume = ob.bids.volume.Sub(orderVolume)
			}

			order.setStatusToFilled()
			break

		} else {
			Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s", order.shortOrderID(), order.Volume(), order.Volume().Sub(marketOrderVolume), marketOrder.shortOrderID(), marketOrder.Volume(), decimal.Zero))

			marketOrders.Remove(marketOrderNode)

			if order.Side() == Buy {
				ob.asks.volume = ob.asks.volume.Sub(marketOrderVolume)
			} else {
				ob.bids.volume = ob.bids.volume.Sub(marketOrderVolume)
			}

			order.setStatusToPartiallyFilled(orderVolume.Sub(marketOrderVolume))
			marketOrder.setStatusToFilled()
		}
	}
}

// func (ob *OrderBook) fillAndRemoveMarketOrder(n *list.Node[*Order]) *Order {
// 	o := n.Value
// 	marketOrders.Remove(n)

// 	if o.Side() == Buy {
// 		ob.marketBuyOrders.Remove(n)
// 	} else {
// 		ob.marketSellOrders.Remove(n)
// 	}
// 	o.setStatusToFilled()
// 	return o
// }

func (ob *OrderBook) fillAndRemoveLimitOrder(n *list.Node[*Order]) *Order {
	o := n.Value

	delete(ob.activeOrders, o.OrderID())

	if o.Side() == Buy {
		ob.bids.Remove(n)
	} else {
		return ob.asks.Remove(n)
	}
	o.setStatusToFilled()
	return o
}

func (ob *OrderBook) PlaceLimitOrder(side Side, clientID uuid.UUID, volume, price decimal.Decimal) (orderID uuid.UUID, err error) {
	if volume.Sign() <= 0 {
		return uuid.Nil, ErrInvalidVolume
	}

	if price.Sign() <= 0 {
		return uuid.Nil, ErrInvalidPrice
	}

	volumeLeft := volume

	o := NewOrder(side, clientID, Limit, price, volume, true)
	Log(fmt.Sprintf("Created limit order: %v", o))

	if o.Side() == Buy && ob.marketSellOrders.Len() > 0 { // there are market orders waiting to be match

		Log(fmt.Sprintf("Limit order matching with market order: %s", o.shortOrderID()))

		ob.matchWithMarketOrders(ob.marketSellOrders, o)

	} else if o.Side() == Sell && ob.marketBuyOrders.Len() > 0 {

		Log(fmt.Sprintf("Limit order matching with market order: %s", o.shortOrderID()))

		ob.matchWithMarketOrders(ob.marketBuyOrders, o)
	}

	if o.Status() == Filled {
		return o.orderID, nil
	}

	var (
		os         *OrderSide
		iter       func() (*OrderQueue, bool)
		comparator func(decimal.Decimal) bool
	)

	if side == Buy {
		iter = ob.asks.MinPriceQueue
		comparator = price.GreaterThanOrEqual
		os = ob.asks
	} else {
		iter = ob.bids.MaxPriceQueue
		comparator = price.LessThanOrEqual
		os = ob.bids
	}

	bestPrice, ok := iter()

	if !ok {
		Log(fmt.Sprintf("No limit orders in the opposite side, initialize order: %s", o.shortOrderID()))
		ob.addLimitOrder(o)
		return o.orderID, nil
	}

	for volumeLeft.Sign() > 0 && os.Len() > 0 && comparator(bestPrice.Price()) {
		bestPrice, _ := iter() // we don't dont have to check ok because we already checked it in the for loop condition with checking orderside size
		volumeLeft = ob.matchAtPriceLevel(bestPrice, o)
	}

	if volumeLeft.Sign() > 0 {
		// the order is not fully filled or didn't find a match in price range, add it to the market order list
		Log(fmt.Sprintf("The order is not fully filled or matched in price range, add it: %s", o.shortOrderID()))
		ob.addLimitOrder(o)
	}

	// o := NewOrder(side, clientID, Limit, price, volume, true)
	// Log(fmt.Sprintf("Created limit order: %v", o))
	// ob.addLimitOrder(o)

	return o.orderID, nil
}

func (ob *OrderBook) addMarketOrder(o *Order) {
	if o.Side() == Buy {
		ob.marketBuyOrders.PushBack(o)
		ob.bids.volume = ob.bids.volume.Add(o.Volume())
	} else {
		ob.marketSellOrders.PushBack(o)
		ob.asks.volume = ob.asks.volume.Add(o.Volume())
	}
}

func (ob *OrderBook) addLimitOrder(o *Order) {
	ob.ordersMu.Lock()
	defer ob.ordersMu.Unlock()

	if o.Side() == Buy {
		ob.activeOrders[o.OrderID()] = ob.bids.Append(o)
	} else {
		ob.activeOrders[o.OrderID()] = ob.asks.Append(o)
	}
}

func (ob *OrderBook) MarketPrice() decimal.Decimal {
	ob.marketPriceMu.RLock()
	defer ob.marketPriceMu.RUnlock()
	return ob.marketPrice
}

func (ob *OrderBook) SetMarketPrice(price decimal.Decimal) {
	ob.marketPriceMu.Lock()
	ob.marketPrice = price
	ob.marketPriceMu.Unlock()

	// release the stop orders that are triggered by the new market price
}
