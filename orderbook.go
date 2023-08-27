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
	orders           map[uuid.UUID]*Order // orderID -> *Order for quick lookup
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
		orders:           map[uuid.UUID]*Order{},
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
	ob.addMarketOrder(o)
	Log(fmt.Sprintf("Created market order: %v", o))
	// fill the order by taking the best limit orders of the opposite Side
	var (
		os           *OrderSide
		iter         func() (*OrderQueue, bool)
		marketOrders *list.List[*Order]
	)
	if side == Buy {
		iter = ob.asks.MinPriceQueue
		os = ob.asks
		marketOrders = ob.marketBuyOrders
	} else {
		iter = ob.bids.MaxPriceQueue
		os = ob.bids
		marketOrders = ob.marketSellOrders
	}

	if os.Len() == 0 {
		// no limit orders in the opposite side, add the order to the market order list
		marketOrders.PushBack(o)
		Log(fmt.Sprintf("No limit orders in the opposite side, add the order to the market order list: %s", o.OrderID()))
		return o.orderID, nil
	}

	var volumeLeft = o.Volume()

	for volumeLeft.Sign() > 0 && os.Len() > 0 { // while the order is not fully filled and the opposite side has more limit orders
		oq, _ := iter()
		volumeLeft = ob.matchAtPriceLevel(os, oq, o)
	}

	if volumeLeft.Sign() > 0 {
		// the order is not fully filled, add it to the market order list
		Log(fmt.Sprintf("The order is not fully filled, add it to the market order list: %s", o.shortOrderID()))
		marketOrders.PushBack(o)
	}
	return o.orderID, nil
}

func (ob *OrderBook) matchAtPriceLevel(os *OrderSide, oq *OrderQueue, o *Order) (volumeLeft decimal.Decimal) {
	volumeLeft = o.Volume()
	Log(fmt.Sprintf("Matching %s at price level %s", o.shortOrderID(), oq.Price()))
	for oq.Len() > 0 && volumeLeft.Sign() > 0 { // while the price level have more limit orders and the order is not fully filled
		bestOrderNode := oq.Head()
		bestOrder := oq.Head().Value
		if volumeLeft.LessThan(bestOrder.Volume()) { // the best order will be partially filled
			matchedVolumeLeft := bestOrder.Volume().Sub(volumeLeft) // update order status. This change should reflect in order queue
			oq.volume = oq.volume.Sub(volumeLeft)
			os.volume = os.volume.Sub(volumeLeft)
			bestOrder.partiallyFillOrder(matchedVolumeLeft)
			o.fillOrder()
			volumeLeft = decimal.Zero
		} else {
			volumeLeft = volumeLeft.Sub(bestOrder.Volume())
			o.partiallyFillOrder(volumeLeft)
			bestOrder.fillOrder()
			oq.Remove(bestOrderNode)
		}
	}
	return
}

func (ob *OrderBook) PlaceLimitOrder(side Side, clientID uuid.UUID, volume, price decimal.Decimal) (orderID uuid.UUID, err error) {
	if volume.Sign() <= 0 {
		return uuid.Nil, ErrInvalidVolume
	}

	if price.Sign() <= 0 {
		return uuid.Nil, ErrInvalidPrice
	}

	o := NewOrder(side, clientID, Limit, price, volume, true)
	Log(fmt.Sprintf("Created limit order: %v", o))
	ob.addLimitOrder(o)

	return o.orderID, nil
}

func (ob *OrderBook) addMarketOrder(o *Order) {
	ob.ordersMu.Lock()
	ob.orders[o.OrderID()] = o
	ob.ordersMu.Unlock()
}

func (ob *OrderBook) addLimitOrder(o *Order) {
	ob.ordersMu.Lock()
	ob.orders[o.OrderID()] = o
	ob.ordersMu.Unlock()

	if o.Side() == Buy {
		ob.bids.Append(o)
	} else {
		ob.asks.Append(o)
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
