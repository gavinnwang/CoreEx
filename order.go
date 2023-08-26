package orderbook

import (
	// "encoding/json"
	// "log"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	side           Side
	orderID        uuid.UUID
	clientID       uuid.UUID
	timestamp      time.Time
	orderType      OrderType
	partialAllowed bool
	status         OrderStatus
	price          decimal.Decimal
	volume         decimal.Decimal
}

func NewOrder (side Side, clientID uuid.UUID, orderType OrderType, price, volume decimal.Decimal, partialAllowed bool) *Order {
	return &Order{
		side:           side,
		orderID:        uuid.New(),
		clientID:       clientID,
		timestamp:      time.Now(),
		orderType:      orderType,
		partialAllowed: partialAllowed,
		status:         Open,
		price:          price,
		volume:         volume,
	}
}

// ID returns orderID field copy
func (o *Order) OrderID() uuid.UUID {
	return o.orderID
}

// Side returns side of the order
func (o *Order) Side() Side {
	return o.side
}

// volume returns volume field copy
func (o *Order) Volume() decimal.Decimal {
	return o.volume
}

// Price returns price field copy
func (o *Order) Price() decimal.Decimal {
	return o.price
}

// Time returns timestamp field copy in Unix
func (o *Order) Time() time.Time {
	return o.timestamp
}

func (o *Order) OrderType() OrderType {
	return o.orderType
}

func (o *Order) ClientID() uuid.UUID {
	return o.clientID
}

// String implements Stringer interface
func (o *Order) String() string {
	return fmt.Sprintf("\n\"%s\":\n\tside: %s\n\ttype: %s\n\tvolume: %s\n\tprice: %s\n\ttime: %s\n", o.OrderID(), o.Side(), o.OrderType(), o.Volume(), o.Price(), o.Time())
}

// // MarshalJSON implements json.Marshaler interface
// func (o *Order) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(
// 		&struct {
// 			S         Side            `json:"side"`
// 			OrderID   string          `json:"orderID"`
// 			Timestamp time.Time       `json:"timestamp"`
// 			Volume    decimal.Decimal `json:"volume"`
// 			Price     decimal.Decimal `json:"price"`
// 		}{
// 			S:         o.Side(),
// 			OrderID:   o.OrderID(),
// 			Timestamp: o.Time(),
// 			Volume:    o.Volume(),
// 			Price:     o.Price(),
// 		},
// 	)
// }

// // UnmarshalJSON implements json.Unmarshaler interface
// func (o *Order) UnmarshalJSON(data []byte) error {
// 	obj := struct {
// 		S         Side            `json:"side"`
// 		OrderID   string          `json:"orderID"`
// 		Timestamp time.Time       `json:"timestamp"`
// 		Volume    decimal.Decimal `json:"volume"`
// 		Price     decimal.Decimal `json:"price"`
// 	}{}

// 	if err := json.Unmarshal(data, &obj); err != nil {
// 		return err
// 	}
// 	orderID, err := uuid.Parse(obj.OrderID)
// 	if err != nil {
// 		log.Fatalf("failed to parse UUID: %v", err)
// 	}
// 	o.orderID = orderID
// 	o.side = obj.S
// 	o.timestamp = obj.Timestamp
// 	o.volume = obj.Volume
// 	o.price = obj.Price
// 	return nil
// }
