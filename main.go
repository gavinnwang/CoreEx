package main

import (
	"encoding/json"
	"net/http"
	"github.com/labstack/echo"
	"github.com/wry0313/crypto-exchange/orderbook"
)

func main() {
	e := echo.New()
	ex := NewExchange()

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)
	e.Start(":3000")
}

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)

type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()
	return &Exchange{
		orderbooks: orderbooks,
	}
}

type PlaceOrderRequest struct {
	Type   OrderType // Limit or Market
	Bid    bool
	Size   float64
	Price  float64
	Market Market
}

// only limit order because market order doesn't have a price
type Order struct {
	Price     float64
	Size      float64
	Bid       bool
	Timestamp int64
}

type OrderbookData struct {
	Asks []*Order
	Bids []*Order
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderbookData := OrderbookData{
		Asks: []*Order{},
		Bids: []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

 return c.JSON(http.StatusOK, orderbookData) 
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	ob := ex.orderbooks[market]
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)
	ob.PlaceLimitOrder(placeOrderData.Price, order)

	return c.JSON(200, map[string]any{"msg": "order placed"})
}
