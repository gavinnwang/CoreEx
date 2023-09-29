package orderbook

import (
	"context"
	"encoding/json"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	list "github/wry-0313/exchange/pkg/dsa/linkedlist"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type Service interface {
	PlaceMarketOrder(side Side, userID ulid.ULID, volume decimal.Decimal) (orderID ulid.ULID, err error)
	PlaceLimitOrder(side Side, userID ulid.ULID, volume, price decimal.Decimal) (orderID ulid.ULID, err error)
	Symbol() string
	NewOrder(side Side, userID ulid.ULID, orderType OrderType, price, volume decimal.Decimal, partialAllowed bool) *Order
	PersistMarketPrice(priceData models.StockPriceHistory) error
	GetMarketPriceHistory() ([]models.StockPriceHistory, error)
	SimulateMarketFluctuations(marketSimulationUlid ulid.ULID)
	Run()
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

	rdb *redis.Client

	prices []decimal.Decimal
}

func NewService(symbol string, obRepo Repository, rdb *redis.Client) Service {
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
		rdb:              rdb,
		prices:           []decimal.Decimal{},
	}
}

func (s *service) Run() {
	loc, _ := time.LoadLocation("America/Chicago")
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:

				var priceData models.StockPriceHistory
				var new bool // new means to create a new candle, otherwise update the last candle
				s.prices = append(s.prices, s.MarketPrice())
				if len(s.prices) >= 5 {
					priceData = models.StockPriceHistory{
						PriceData:  getPriceDataFromPriceSlice(s.prices),
						BidVolume:  s.bids.Volume().InexactFloat64(),
						AskVolume:  s.asks.Volume().InexactFloat64(),
						RecordedAt: time.Now().In(loc),
					}
					new = true

					err := s.PersistMarketPrice(priceData)
					if err != nil {
						log.Printf("Could not persist market price: %v", err)
					}
					s.prices = []decimal.Decimal{}
					s.asks.ResetVolume()
					s.bids.ResetVolume()
				} else {
					priceData = models.StockPriceHistory{
						PriceData:  getPriceDataFromPriceSlice(s.prices),
						BidVolume:  s.bids.Volume().InexactFloat64(),
						AskVolume:  s.asks.Volume().InexactFloat64(),
						RecordedAt: time.Now().In(loc),
					}

					new = false
				}

				s.publishPrice(priceData, new)

			}
		}
	}()
}

func getPriceDataFromPriceSlice(prices []decimal.Decimal) models.PriceData {

	return models.PriceData{
		Open:  prices[0],
		Close: prices[max(0, len(prices)-1)],
		High:  getMax(prices),
		Low:   getMin(prices),
	}

}

func getMin(prices []decimal.Decimal) decimal.Decimal {
	min := prices[0]
	for _, price := range prices {
		if price.LessThan(min) {
			min = price
		}
	}
	return min
}

func getMax(prices []decimal.Decimal) decimal.Decimal {
	max := prices[0]
	for _, price := range prices {
		if price.GreaterThan(max) {
			max = price
		}
	}
	return max
}

func (s *service) publishPrice(priceData models.StockPriceHistory, new bool) {

	symbolMarketInfo := SymbolInfoResponse{
		Symbol:  s.symbol,
		Price:   s.MarketPrice().InexactFloat64(),
		BestBid: s.BestBid().InexactFloat64(),
		BestAsk: s.BestAsk().InexactFloat64(),
		// AskVolume: s.asks.Volume().InexactFloat64(),
		// BidVolume: s.bids.Volume().InexactFloat64(),
		CandleData: CandleData{
			StockPriceHistory: priceData,
			NewCandle:         new,
		},
	}
	// log.Printf("Publishing market info: %v to redis channel %s\n", symbolMarketInfo, s.symbol)

	pubMsg := SymbolInfoPubMsg{
		RedisPubMsgBase: RedisPubMsgBase{
			Event:   EventStreamSymbolInfo,
			Success: true,
		},
		Result: symbolMarketInfo,
	}
	pubMsgBytes, err := json.Marshal(pubMsg)
	if err != nil {
		log.Printf("Service: failed to marshal market info into JSON: %v", s.symbol)
		return
	}

	s.rdb.Publish(context.Background(), s.symbol, pubMsgBytes)
}

// Define our sine wave parameters
type SineWave struct {
	frequency float64
	amplitude float64
	phase     float64
}

var sineWaves = []SineWave{

	{frequency: 0.1, amplitude: 2, phase: 0.5},
	{frequency: 0.3, amplitude: 1, phase: 0.2},

	{frequency: 0.4, amplitude: 10, phase: 0.3},
	{frequency: 0.6, amplitude: 5, phase: 0.1},
	{frequency: 0.8, amplitude: 2, phase: 0.5},

	{frequency: 0.1, amplitude: 3, phase: 0},
	{frequency: 0.3, amplitude: 1.5, phase: 1},
	{frequency: 0.5, amplitude: 1, phase: 0.5},
	{frequency: 0.7, amplitude: 0.5, phase: 0.2},
}

func calculateSuperimposedSine(t float64) float64 {
	totalValue := 0.0
	for _, wave := range sineWaves {
		totalValue += wave.amplitude * math.Sin(wave.frequency*t+wave.phase)
	}
	// log.Printf("Total value: %f\	 totalValue);
	return totalValue / 10
}

func (s *service) SimulateMarketFluctuations(marketSimulationUlid ulid.ULID) {

	// s.PlaceLimitOrder(Buy, marketSimulationUlid, decimal.NewFromInt(100), decimal.NewFromInt(145))
	// s.PlaceLimitOrder(Sell, marketSimulationUlid, decimal.NewFromInt(100), decimal.NewFromInt(150))

	s.SetMarketPrice(decimal.NewFromFloat(150))

	// time.Sleep(3 * time.Second)
	t := 0.0

	go func() { // limit buy order
		for {
			// log.Printf("Best ask: %s", s.BestAsk())
			priceFluctuation := calculateSuperimposedSine(t)
			price := (s.MarketPrice()).Add(decimal.NewFromFloat(3)).Sub(decimal.NewFromFloat(rand.Float64() * 5)).Add(decimal.NewFromFloat(priceFluctuation)).Round(2)
			volume := decimal.NewFromFloat(rand.Float64() * 20).Add(decimal.NewFromInt(100)).Round(2)
			_, err := s.PlaceLimitOrder(Buy, marketSimulationUlid, volume, price)
			if err != nil {
				// log.Printf("Failed to place limit order: %v", err)
			}
			time.Sleep(50 * time.Millisecond)
			t += 0.03
		}
	}()

	t1 := 0.0
	
	go func() {
		for {
			// log.Printf("Best bid: %s", s.BestBid())
			priceFluctuation := calculateSuperimposedSine(t1)
			price := (s.MarketPrice()).Sub(decimal.NewFromFloat(3)).Add(decimal.NewFromFloat(rand.Float64() * 5)).Add(decimal.NewFromFloat(priceFluctuation)).Round(2)
			volume := decimal.NewFromFloat(rand.Float64() * 20).Add(decimal.NewFromInt(100)).Round(2)
			_, err := s.PlaceLimitOrder(Sell, marketSimulationUlid, volume, price)
			if err != nil {
				// log.Printf("Failed to place limit order: %v", err)
			}
			time.Sleep(50 * time.Millisecond)
			t1 += 0.03
		}
	}()

	// // time.Sleep(3 * time.Second)
	// t2 := 0.0

	// go func() { // limit buy order
	// 	for {
	// 		// log.Printf("Best ask: %s", s.BestAsk())
	// 		priceFluctuation := calculateSuperimposedSine(t2)
	// 		price := (s.MarketPrice()).Add(decimal.NewFromFloat(3)).Sub(decimal.NewFromFloat(rand.Float64() * 5)).Add(decimal.NewFromFloat(priceFluctuation)).Round(2)
	// 		volume := decimal.NewFromFloat(rand.Float64() * 20).Add(decimal.NewFromInt(100)).Round(2)
	// 		_, err := s.PlaceLimitOrder(Buy, marketSimulationUlid, volume, price)
	// 		if err != nil {
	// 			// log.Printf("Failed to place limit order: %v", err)
	// 		}
	// 		time.Sleep(50 * time.Millisecond)
	// 		t2 += 0.03
	// 	}
	// }()

	// t3 := 0.0

	// go func() {
	// 	for {
	// 		// log.Printf("Best bid: %s", s.BestBid())
	// 		priceFluctuation := calculateSuperimposedSine(t3)
	// 		price := (s.MarketPrice()).Sub(decimal.NewFromFloat(3)).Add(decimal.NewFromFloat(rand.Float64() * 5)).Add(decimal.NewFromFloat(priceFluctuation)).Round(2)
	// 		volume := decimal.NewFromFloat(rand.Float64() * 20).Add(decimal.NewFromInt(100)).Round(2)
	// 		_, err := s.PlaceLimitOrder(Sell, marketSimulationUlid, volume, price)
	// 		if err != nil {
	// 			// log.Printf("Failed to place limit order: %v", err)
	// 		}
	// 		time.Sleep(50 * time.Millisecond)
	// 		t3 += 0.03
	// 	}
	// }()

	// go func() { // limit buy order
	// 	for {
	// 		log.Printf("Best ask: %s", s.BestAsk())
	// 		price := decimal.NewFromFloat(rand.Float64() * 5).Add(s.BestAsk()).Sub(decimal.NewFromInt(1)).Round(2)
	// 		volume := decimal.NewFromFloat(rand.Float64() * 20).Add(decimal.NewFromInt(100)).Round(2)
	// 		_, err := s.PlaceLimitOrder(Buy, marketSimulationUlid, volume, price)
	// 		if err != nil {
	// 			log.Printf("Failed to place limit order: %v", err)
	// 		}
	// 		time.Sleep(100 * time.Millisecond)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		log.Printf("Best bid: %s", s.BestBid())
	// 		price := decimal.NewFromFloat(rand.Float64() * 5).Add(s.BestBid()).Add(decimal.NewFromInt(1)).Round(2)
	// 		volume := decimal.NewFromFloat(rand.Float64() * 20).Add(decimal.NewFromInt(100)).Round(2)
	// 		_, err := s.PlaceLimitOrder(Sell, marketSimulationUlid, volume, price)
	// 		if err != nil {
	// 			log.Printf("Failed to place limit order: %v", err)
	// 		}
	// 		time.Sleep(100 * time.Millisecond)
	// 	}
	// }()

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

	// os.AddVolumeBy(volume)

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

			// if o.Side() == Buy {
			// 	s.asks.SubVolumeBy(volumeLeft)
			// } else {
			// 	s.bids.SubVolumeBy(volumeLeft)
			// }
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

			// if order.Side() == Buy {
			// 	s.asks.SubVolumeBy(orderVolume)
			// } else {
			// 	s.bids.SubVolumeBy(orderVolume)
			// }

			break

		} else {
			// Log(fmt.Sprintf("%s: %s -> %s | %s: %s -> %s\n", order.shortOrderID(), order.Volume(), order.Volume().Sub(marketOrderVolume), marketOrder.shortOrderID(), marketOrder.Volume(), decimal.Zero))

			marketOrders.Remove(marketOrderNode)

			// if order.Side() == Buy {
			// 	s.asks.SubVolumeBy(marketOrderVolume)
			// } else {
			// 	s.bids.SubVolumeBy(marketOrderVolume)
			// }

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

	// os.AddVolumeBy(volume);

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
	// log.Printf("Persisting market price: %v", priceData)
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

// func (s *service) AskVolume() decimal.Decimal {
// 	return s.asks.Volume()
// }

// func (s *service) BidVolume() decimal.Decimal {
// 	return s.bids.Volume()
// }

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
