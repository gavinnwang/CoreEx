package orderbook

import (
	list "github/wry-0313/exchange/linkedlist"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	orders       map[uuid.UUID]*Order // orderID -> *Order
	currentPrice decimal.Decimal
	marketOrders *list.List[*Order]
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
  if volume.Sign() <= 0 {
    return 
  }
	// create the order and add it to all orders
	o := NewOrder(side, clientID, Market, decimal.Zero, volume, true)
	ob.orders[o.OrderID()] = o
	// fill the order by taking the best limit orders of the opposite Side
	var oq *OrderQueue
	var ok bool
	if side == Buy {
		oq, ok = ob.asks.MinPriceQueue()
	} else {
		oq, ok = ob.asks.MaxPriceQueue()
	}

	if !ok {
		ob.marketOrders.PushBack(o)
    return 
	}

  for volume.Sign() > 0 && 
}

func (ob *OrderBook) matchAtPriceLevel(oq *OrderQueue, o *Order) (volumeLeft decimal.Decimal) {
	volumeLeft = o.Volume()
	// while the price level have more limit orders
	for oq.Len() > 0 && volumeLeft.Sign() > 0 {
		bestOrderNode := oq.Head()
		bestOrder := oq.Head().Value
		if volumeLeft.LessThan(bestOrder.Volume()) {
			// the best order will be partially filled
			matchedVolumeLeft := bestOrder.Volume().Sub(volumeLeft)
			// update order status. This change should reflect in order queue
			bestOrder.partiallyFillOrder(matchedVolumeLeft)
			o.fillOrder()
		} else {
			volumeLeft = volumeLeft.Sub(bestOrder.Volume())
			o.partiallyFillOrder(volumeLeft)
			bestOrder.fillOrder()
			oq.Remove(bestOrderNode)
		}
	}
	return
}
