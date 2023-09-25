package orderbook

import (
	"fmt"
	"github/wry-0313/exchange/internal/models"
	list "github/wry-0313/exchange/pkg/dsa/linkedlist"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type Service interface {
	PlaceMarketOrder(side Side, userID ulid.ULID, volume decimal.Decimal) (orderID ulid.ULID, err error)
	PlaceLimitOrder(side Side, userID ulid.ULID, volume, price decimal.Decimal) (orderID ulid.ULID, err error)
	AskVolume() decimal.Decimal
	BidVolume() decimal.Decimal
	BestBid() decimal.Decimal
	BestAsk() decimal.Decimal
	MarketPrice() decimal.Decimal
	Symbol() string
	NewOrder(side Side, userID ulid.ULID, orderType OrderType, price, volume decimal.Decimal, partialAllowed bool) *Order
	PersistMarketPrice(priceData models.StockPriceHistory) error
	GetMarketPriceHistory() ([]models.StockPriceHistory, error)
	RunMarketPriceHistoryPersistance()
	SimulateMarketFluctuations(marketSimulationUlid ulid.ULID)
}

type service struct {
	symbol       string
	activeOrders map[ulid.ULID]*list.Node[*Order] // orderID -> *Order for quick acctions such as update or cancel
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

	obRepo Repository

	prices []decimal.Decimal
}

func NewService(symbol string, obRepo Repository) Service {
	err := InitializeLogService("orderbook_log.txt")
	if err != nil {
		log.Fatalf("Could not initialize log service: %v", err)
	}

	stock := models.Stock{
		Symbol: symbol,
	}

	err = obRepo.CreateStock(stock)
	if err != nil {
		log.Fatalf("Could not create stock: %v", err)
	}

	return &service{
		symbol:           symbol,
		activeOrders:     map[ulid.ULID]*list.Node[*Order]{},
		bids:             NewOrderSide(),
		asks:             NewOrderSide(),
		marketBuyOrders:  list.New[*Order](),
		marketSellOrders: list.New[*Order](),
		marketPrice:      decimal.Zero,
		obRepo:           obRepo,
		prices:           []decimal.Decimal{},
	}
}

func (s *service) SimulateMarketFluctuations(marketSimulationUlid ulid.ULID) {
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			side := Buy
			if rand.Intn(2) == 0 {
				side = Sell
			}

			volume := decimal.NewFromFloat(rand.Float64() * 10).Round(2)
			price := decimal.NewFromFloat(rand.Float64() * 5000).Round(2)
			log.Printf("Simulating market fluctuations: %v %v %v\n", side, volume, price)
			if rand.Intn(2) == 0 {
				_, _ = s.PlaceLimitOrder(side, marketSimulationUlid, volume, price)
			} else {
				_, _ = s.PlaceLimitOrder(side, marketSimulationUlid, volume, price)
			}
		}
	}()
}

func (s *service) RunMarketPriceHistoryPersistance() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// log.Printf("price vector %v\n", s.prices)
				if len(s.prices) >= 5 {
					open := s.prices[0]
					close := s.prices[len(s.prices)-1]
					sort.Slice(s.prices, func(i, j int) bool {
						return s.prices[i].LessThan(s.prices[j])
					})
					priceData := models.StockPriceHistory{
						Open:       open,
						Close:      close,
						High:       s.prices[len(s.prices)-1],
						Low:        s.prices[0],
						Volume:     s.AskVolume().Add(s.BidVolume()),
						RecordedAt: time.Now().Unix(),
					}
					err := s.PersistMarketPrice(priceData)
					if err != nil {
						log.Printf("Could not persist market price: %v", err)
					}
					s.prices = []decimal.Decimal{}
				}
				s.prices = append(s.prices, s.MarketPrice())
			}
		}
	}()
}

func (s *service) GetMarketPriceHistory() ([]models.StockPriceHistory, error) {
	return s.obRepo.GetEntireMarketPriceHistory(s.symbol)
}

func (s *service) PlaceMarketOrder(side Side, userID ulid.ULID, volume decimal.Decimal) (orderID ulid.ULID, err error) {
	if volume.Sign() <= 0 {
		return ulid.ULID{}, ErrInvalidVolume
	}

	if side == Invalid {
		return ulid.ULID{}, ErrInvalidSide
	}

	o := s.NewOrder(side, userID, Market, decimal.Zero, volume, true)

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

	logService.logger.Println(fmt.Sprintf("Matching %s at price level %s\n", o.shortOrderID(), oq.Price()))

	s.SetMarketPrice(oq.Price())

	for oq.Len() > 0 && volumeLeft.Sign() > 0 { // while the price level have more limit orders and the order is not fully filled

		bestOrderNode := oq.Head()
		bestOrder := oq.Head().Value

		bestOrderVolume := bestOrder.Volume()

		logService.logger.Println(fmt.Sprintf("Matching %s with %s", o.shortOrderID(), bestOrder.shortOrderID()))

		if volumeLeft.LessThan(bestOrderVolume) { // the best order will be partially filled

			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", o.shortOrderID(), o.Volume(), o.Volume().Sub(volumeLeft), bestOrder.shortOrderID(), bestOrder.Volume(), bestOrder.Volume().Sub(volumeLeft)))
			// matchedVolumeLeft := bestOrderVolume.Sub(volumeLeft) // update order status. This change should reflect in order queue
			oq.SetVolume(oq.Volume().Sub(volumeLeft))

			if o.Side() == Buy {
				s.asks.SubVolumeBy(volumeLeft)
			} else {
				s.bids.SubVolumeBy(volumeLeft)
			}
			s.fillOrder(bestOrder, volumeLeft, oq.Price())
			s.fillOrder(o, volumeLeft, oq.Price()) // completely filled

			volumeLeft = decimal.Zero

		} else { // the best order will be completely filled
			volumeLeft = volumeLeft.Sub(bestOrderVolume)
			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", o.shortOrderID(), o.Volume(), o.Volume().Sub(bestOrder.Volume()), bestOrder.shortOrderID(), bestOrder.Volume(), decimal.Zero))
			s.fillAndRemoveLimitOrder(bestOrderNode, bestOrderVolume, oq.Price())
			s.fillOrder(o, bestOrderVolume, oq.Price())
		}
	}
	return
}

// func (s *service) processTransaction()

func (s *service) matchWithMarketOrders(marketOrders *list.List[*Order], order *Order) {

	for marketOrders.Len() > 0 {

		s.SetMarketPrice(order.Price())

		marketOrderNode := marketOrders.Front()
		marketOrder := marketOrderNode.Value
		marketOrderVolume := marketOrder.Volume()
		orderVolume := order.Volume()

		if orderVolume.LessThan(marketOrderVolume) { // the market order will be completely filled

			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", order.shortOrderID(), order.Volume(), decimal.Zero, marketOrder.shortOrderID(), marketOrder.Volume(), marketOrder.Volume().Sub(orderVolume)))

			s.fillOrder(marketOrder, orderVolume, order.Price())
			s.fillOrder(order, orderVolume, order.Price())

			if order.Side() == Buy {
				s.asks.SubVolumeBy(orderVolume)
			} else {
				s.bids.SubVolumeBy(orderVolume)
			}

			break

		} else {
			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", order.shortOrderID(), order.Volume(), order.Volume().Sub(marketOrderVolume), marketOrder.shortOrderID(), marketOrder.Volume(), decimal.Zero))

			marketOrders.Remove(marketOrderNode)

			if order.Side() == Buy {
				s.asks.SubVolumeBy(marketOrderVolume)
			} else {
				s.bids.SubVolumeBy(marketOrderVolume)
			}

			s.fillOrder(order, marketOrderVolume, order.Price())
			s.fillOrder(marketOrder, marketOrderVolume, order.Price())
		}
	}
}

func (s *service) fillAndRemoveLimitOrder(n *list.Node[*Order], filledVolume, filledAt decimal.Decimal) *Order {
	o := n.Value

	s.ordersMu.Lock()
	delete(s.activeOrders, o.OrderID())
	s.ordersMu.Unlock()

	if o.Side() == Buy {
		s.bids.Remove(n)
	} else {
		s.asks.Remove(n)
	}
	s.fillOrder(o, filledVolume, filledAt)
	return o
}

func (s *service) PlaceLimitOrder(side Side, userID ulid.ULID, volume, price decimal.Decimal) (orderID ulid.ULID, err error) {
	if volume.Sign() <= 0 {
		return ulid.ULID{}, ErrInvalidVolume
	}

	if price.Sign() <= 0 {
		return ulid.ULID{}, ErrInvalidPrice
	}

	if side == Invalid {
		return ulid.ULID{}, ErrInvalidSide
	}

	volumeLeft := volume

	o := s.NewOrder(side, userID, Limit, price, volume, true)

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
		logService.logger.Println(fmt.Sprintf("No limit orders in the opposite side, initialize order: %s", o.shortOrderID()))
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

func (s *service) PersistMarketPrice(priceData models.StockPriceHistory) error {
	log.Printf("Persisting market price: %v", priceData)
	// persist market price to db
	err := s.obRepo.CreateMarketPriceHistory(s.symbol, priceData)
	if err != nil {
		return fmt.Errorf("Service: failed to persist market price history: %w", err)
	}
	return nil
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
	logService.logger.Println(fmt.Sprintf("Set market price: %s", price))
	s.marketPrice = price
	s.marketPriceMu.Unlock()

	// release the stop orders that are triggered by the new market price
}

func (s *service) Symbol() string {
	return s.symbol
}
