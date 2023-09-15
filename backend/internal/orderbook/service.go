package orderbook

import (
	"fmt"
	list "github/wry-0313/exchange/pkg/dsa/linkedlist"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	PlaceMarketOrder(side Side, userID uuid.UUID, volume decimal.Decimal) (orderID uuid.UUID, err error)
	PlaceLimitOrder(side Side, userID uuid.UUID, volume, price decimal.Decimal) (orderID uuid.UUID, err error)
	AskVolume() decimal.Decimal
	BidVolume() decimal.Decimal
	BestBid() decimal.Decimal
	BestAsk() decimal.Decimal
	MarketPrice() decimal.Decimal
	Symbol() string
}

type service struct {
	symbol       string
	activeOrders map[uuid.UUID]*list.Node[*Order] // orderID -> *Order for quick acctions such as update or cancel
	ordersMu     sync.RWMutex

	marketPrice   decimal.Decimal
	marketPriceMu sync.RWMutex

	marketBuyOrders *list.List[*Order] // partially filled market buy orders
	marketBuyMu     sync.Mutex

	marketSellOrders *list.List[*Order] // partially filled market sell orders
	marketSellMu     sync.Mutex

	asks *OrderSide // limit sell orders
	bids *OrderSide // limit buy orders

	sortedOrdersMu sync.RWMutex
}

func NewService(symbol string) Service {
	err := InitializeLogService("orderbook_log.txt")
	if err != nil {
		log.Fatalf("Could not initialize log service: %v", err)
	}
	return &service{
		symbol:           symbol,
		activeOrders:     map[uuid.UUID]*list.Node[*Order]{},
		bids:             NewOrderSide(),
		asks:             NewOrderSide(),
		marketBuyOrders:  list.New[*Order](),
		marketSellOrders: list.New[*Order](),
		marketPrice:      decimal.Zero,
	}
}

func (s *service) PlaceMarketOrder(side Side, userID uuid.UUID, volume decimal.Decimal) (orderID uuid.UUID, err error) {
	if volume.Sign() <= 0 {
		return uuid.Nil, ErrInvalidVolume
	}

	if side == Invalid {
		return uuid.Nil, ErrInvalidSide
	}

	o := NewOrder(side, userID, Market, decimal.Zero, volume, true)
	Log(fmt.Sprintf("Created market order: %v", o))

	var (
		os   *OrderSide
		iter func() (*OrderQueue, bool)
	)
	if side == Buy {
		iter = s.asks.MinPriceQueue
		os = s.asks
	} else {
		iter = s.bids.MaxPriceQueue
		os = s.bids
	}

	if os.Len() == 0 {
		// no limit orders in the opposite side, add the order to the market order list
		s.addMarketOrder(o)
		return o.orderID, nil
	}

	var volumeLeft = o.Volume()

	s.sortedOrdersMu.Lock()
	for volumeLeft.Sign() > 0 && os.Len() > 0 { // while the order is not fully filled and the opposite side has more limit orders)

		oq, found := iter()

		if !found {
			continue
		}

		// Log(fmt.Sprintf("os string: %v, oq: %v  volumeleft: %s\n", os, oq, volumeLeft))

		volumeLeft = s.matchAtPriceLevel(oq, o)
	}
	s.sortedOrdersMu.Unlock()

	if volumeLeft.Sign() > 0 {
		// the order is not fully filled, add it to the market order list
		s.addMarketOrder(o)
	} else {
		// the order is fully filled
	}
	return o.orderID, nil
}

func (s *service) matchAtPriceLevel(oq *OrderQueue, o *Order) (volumeLeft decimal.Decimal) {
	volumeLeft = o.Volume()

	Log(fmt.Sprintf("Matching %s at price level %s\n", o.shortOrderID(), oq.Price()))

	s.SetMarketPrice(oq.Price())

	for oq.Len() > 0 && volumeLeft.Sign() > 0 { // while the price level have more limit orders and the order is not fully filled

		bestOrderNode := oq.Head()
		bestOrder := oq.Head().Value

		Log(fmt.Sprintf("Matching %s with %s", o.shortOrderID(), bestOrder.shortOrderID()))

		if volumeLeft.LessThan(bestOrder.Volume()) { // the best order will be partially filled

			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", o.shortOrderID(), o.Volume(), o.Volume().Sub(volumeLeft), bestOrder.shortOrderID(), bestOrder.Volume(), bestOrder.Volume().Sub(volumeLeft)))

			matchedVolumeLeft := bestOrder.Volume().Sub(volumeLeft) // update order status. This change should reflect in order queue
			oq.SetVolume(oq.Volume().Sub(volumeLeft))

			if o.Side() == Buy {
				s.asks.SubVolumeBy(volumeLeft)
			} else {
				s.bids.SubVolumeBy(volumeLeft)
			}

			bestOrder.setStatusToPartiallyFilled(matchedVolumeLeft)
			o.setStatusToFilled()

			volumeLeft = decimal.Zero

		} else { // the best order will be completely filled
			volumeLeft = volumeLeft.Sub(bestOrder.Volume())

			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", o.shortOrderID(), o.Volume(), o.Volume().Sub(bestOrder.Volume()), bestOrder.shortOrderID(), bestOrder.Volume(), decimal.Zero))

			s.fillAndRemoveLimitOrder(bestOrderNode)
			o.setStatusToPartiallyFilled(volumeLeft)
		}
	}
	return
}

func (s *service) matchWithMarketOrders(marketOrders *list.List[*Order], order *Order) {

	for marketOrders.Len() > 0 {

		s.SetMarketPrice(order.Price())

		marketOrderNode := marketOrders.Front()
		marketOrder := marketOrderNode.Value
		marketOrderVolume := marketOrder.Volume()
		orderVolume := order.Volume()

		if orderVolume.LessThan(marketOrderVolume) { // the market order will be completely filled

			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", order.shortOrderID(), order.Volume(), decimal.Zero, marketOrder.shortOrderID(), marketOrder.Volume(), marketOrder.Volume().Sub(orderVolume)))

			marketOrder.setStatusToPartiallyFilled(marketOrderVolume.Sub(orderVolume))

			if order.Side() == Buy {
				s.asks.SubVolumeBy(orderVolume)
			} else {
				s.bids.SubVolumeBy(orderVolume)
			}

			order.setStatusToFilled()
			break

		} else {
			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", order.shortOrderID(), order.Volume(), order.Volume().Sub(marketOrderVolume), marketOrder.shortOrderID(), marketOrder.Volume(), decimal.Zero))

			marketOrders.Remove(marketOrderNode)

			if order.Side() == Buy {
				s.asks.SubVolumeBy(marketOrderVolume)
			} else {
				s.bids.SubVolumeBy(marketOrderVolume)
			}

			order.setStatusToPartiallyFilled(orderVolume.Sub(marketOrderVolume))
			marketOrder.setStatusToFilled()
		}
	}
}

func (s *service) fillAndRemoveLimitOrder(n *list.Node[*Order]) *Order {
	o := n.Value

	s.ordersMu.Lock()
	delete(s.activeOrders, o.OrderID())
	s.ordersMu.Unlock()

	if o.Side() == Buy {
		s.bids.Remove(n)
	} else {
		return s.asks.Remove(n)
	}
	o.setStatusToFilled()
	return o
}

func (s *service) PlaceLimitOrder(side Side, userID uuid.UUID, volume, price decimal.Decimal) (orderID uuid.UUID, err error) {
	if volume.Sign() <= 0 {
		return uuid.Nil, ErrInvalidVolume
	}

	if price.Sign() <= 0 {
		return uuid.Nil, ErrInvalidPrice
	}

	if side == Invalid {
		return uuid.Nil, ErrInvalidSide
	}

	volumeLeft := volume

	o := NewOrder(side, userID, Limit, price, volume, true)
	Log(fmt.Sprintf("Created limit order: %v", o))

	if o.Side() == Buy { // there are market orders waiting to be match

		s.marketSellMu.Lock() // Lock the mutex
		if s.marketSellOrders.Len() > 0 {
			// Log(fmt.Sprintf("Limit order matching with market order: %s", o.shortOrderID()))
			s.matchWithMarketOrders(s.marketSellOrders, o)
		}
		s.marketSellMu.Unlock() // Unlock the mutex

	} else if o.Side() == Sell {

		s.marketBuyMu.Lock() // Lock the mutex
		if s.marketBuyOrders.Len() > 0 {
			// Log(fmt.Sprintf("Limit order matching with market order: %s", o.shortOrderID()))
			s.matchWithMarketOrders(s.marketBuyOrders, o)
		}
		s.marketBuyMu.Unlock() // Unlock the mutex
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
		iter = s.asks.MinPriceQueue
		comparator = price.GreaterThanOrEqual
		os = s.asks
	} else {
		iter = s.bids.MaxPriceQueue
		comparator = price.LessThanOrEqual
		os = s.bids
	}

	s.sortedOrdersMu.Lock()
	bestPrice, ok := iter()

	if !ok {
		Log(fmt.Sprintf("No limit orders in the opposite side, initialize order: %s", o.shortOrderID()))
		s.addLimitOrder(o)
		s.sortedOrdersMu.Unlock()
		return o.orderID, nil
	}

	for volumeLeft.Sign() > 0 && os.Len() > 0 && comparator(bestPrice.Price()) {
		bestPrice, _ := iter() // we don't dont have to check ok because we already checked it in the for loop condition with checking orderside size
		volumeLeft = s.matchAtPriceLevel(bestPrice, o)
	}

	if volumeLeft.Sign() > 0 {
		// the order is not fully filled or didn't find a match in price range, add it to the market order list
		// Log(fmt.Sprintf("The order is not fully filled or matched in price range, add it: %s", o.shortOrderID()))
		s.addLimitOrder(o)
	}
	s.sortedOrdersMu.Unlock()

	return o.orderID, nil
}

func (s *service) addMarketOrder(o *Order) {
	if o.Side() == Buy {
		s.marketBuyMu.Lock() // Lock the mutex
		s.marketBuyOrders.PushBack(o)
		s.marketBuyMu.Unlock() // Unlock the mutex
		s.bids.AddVolumeBy(o.Volume())
	} else {
		s.marketSellMu.Lock() // Lock the mutex
		s.marketSellOrders.PushBack(o)
		s.marketSellMu.Unlock() // Unlock the mutex
		s.asks.AddVolumeBy(o.Volume())
	}
}

func (s *service) addLimitOrder(o *Order) {
	if o.Side() == Buy {
		n := s.bids.Append(o)

		s.ordersMu.Lock()
		s.activeOrders[o.OrderID()] = n
		s.ordersMu.Unlock()
	} else {
		n := s.asks.Append(o)

		s.ordersMu.Lock()
		s.activeOrders[o.OrderID()] = n
		s.ordersMu.Unlock()
	}
}

func (s *service) AskVolume() decimal.Decimal {
	return s.asks.Volume()
}

func (s *service) BidVolume() decimal.Decimal {
	return s.bids.Volume()
}

func (s *service) BestBid() decimal.Decimal {
	s.sortedOrdersMu.RLock()
	defer s.sortedOrdersMu.RUnlock()
	oq, found := s.bids.MaxPriceQueue()
	if !found {
		return decimal.Zero
	}
	return oq.Price()
}

func (s *service) BestAsk() decimal.Decimal {
	s.sortedOrdersMu.RLock()
	defer s.sortedOrdersMu.RUnlock()
	oq, found := s.asks.MinPriceQueue()
	if !found {
		return decimal.Zero
	}
	return oq.Price()
}

func (s *service) MarketPrice() decimal.Decimal {
	s.marketPriceMu.RLock()
	defer s.marketPriceMu.RUnlock()
	return s.marketPrice
}

func (s *service) SetMarketPrice(price decimal.Decimal) {
	s.marketPriceMu.Lock()
	Log(fmt.Sprintf("Set market price: %s", price))
	s.marketPrice = price
	s.marketPriceMu.Unlock()

	// release the stop orders that are triggered by the new market price
}

func (s *service) Symbol() string {
	return s.symbol
}
