package orderbook

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	side Side
	id   string

	accountID     string
	timestamp int64 // unix.nano for bette perforamnce
	price     decimal.Decimal
	quantity  decimal.Decimal
	isMaker   bool
}

type OrderUpdate struct {
	side      Side
	id        string
	accountID string
	size      decimal.Decimal
	price     decimal.Decimal
}

// ID returns orderID field copy
func (o *Order) ID() string {
	return o.id
}

// Side returns side of the order
func (o *Order) Side() Side {
	return o.side
}

// Quantity returns quantity field copy
func (o *Order) Quantity() decimal.Decimal {
	return o.quantity
}

// Price returns price field copy
func (o *Order) Price() decimal.Decimal {
	return o.price
}

// Time returns timestamp field copy in Unix
func (o *Order) Time() int64 {
	return o.timestamp
}

// Time returns the timestamp field copy in human-readable format
func (o *Order) TimeString() string {
	seconds := o.timestamp / 1e9
	nanos := o.timestamp % 1e9
	t := time.Unix(seconds, nanos)
	return t.Format(time.RFC3339)
}

// String implements Stringer interface
func (o *Order) String() string {
	return fmt.Sprintf("\n\"%s\":\n\tside: %s\n\tquantity: %s\n\tprice: %s\n\ttime: %s\n", o.ID(), o.Side(), o.Quantity(), o.Price(), o.TimeString())
}

// MarshalJSON implements json.Marshaler interface
func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			S         Side            `json:"side"`
			ID        string          `json:"id"`
			Timestamp int64           `json:"timestamp"`
			Quantity  decimal.Decimal `json:"quantity"`
			Price     decimal.Decimal `json:"price"`
		}{
			S:         o.Side(),
			ID:        o.ID(),
			Timestamp: o.Time(),
			Quantity:  o.Quantity(),
			Price:     o.Price(),
		},
	)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (o *Order) UnmarshalJSON(data []byte) error {
	obj := struct {
		S         Side            `json:"side"`
		ID        string          `json:"id"`
		Timestamp int64           `json:"timestamp"`
		Quantity  decimal.Decimal `json:"quantity"`
		Price     decimal.Decimal `json:"price"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	o.side = obj.S
	o.id = obj.ID
	o.timestamp = obj.Timestamp
	o.quantity = obj.Quantity
	o.price = obj.Price
	return nil
}
