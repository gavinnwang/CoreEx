package orderbook

import (
	// "encoding/json"
	// "log"
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type Order struct {
	side      Side
	orderID   ulid.ULID
	userID    ulid.ULID
	orderType OrderType
	status    OrderStatus
	price     decimal.Decimal
	volume    decimal.Decimal
	createdAt time.Time
	volumeMu  sync.RWMutex
}

func NewOrder(side Side, userID ulid.ULID, orderType OrderType, price, volume decimal.Decimal, partialAllowed bool) *Order {
	return &Order{
		side:      side,
		orderID:   ulid.Make(),
		userID:    userID,
		orderType: orderType,
		status:    Open,
		price:     price,
		volume:    volume,
		createdAt: time.Now(),
	}
}

// ID returns orderID field copy
func (o *Order) OrderID() ulid.ULID {
	return o.orderID
}

// shortOrderID returns first 4 characters of orderID (for debugging purposes)
func (o *Order) shortOrderID() string {
	return o.orderID.String()[:6]
}

// Status returns status field copy
func (o *Order) Status() OrderStatus {
	return o.status
}

// Side returns side field copy
func (o *Order) Side() Side {
	return o.side
}

// volume returns volume field copy
func (o *Order) Volume() decimal.Decimal {
	o.volumeMu.RLock()
	defer o.volumeMu.RUnlock()
	return o.volume
}

// Price returns price field copy
func (o *Order) Price() decimal.Decimal {
	return o.price
}

func (o *Order) OrderType() OrderType {
	return o.orderType
}

func (o *Order) UserID() ulid.ULID {
	return o.userID
}

func (o *Order) setStatusToPartiallyFilled(remaining decimal.Decimal) {
	o.volumeMu.Lock()
	o.volume = remaining
	o.volumeMu.Unlock()
	o.status = PartiallyFilled
	// logMsg := fmt.Sprintf("Order partially filled. Remaining volume: %s", remaining)
	// o.AppendLog(logMsg)
}

func (o *Order) setStatusToFilled() {
	o.volumeMu.Lock()
	o.volume = decimal.Zero
	o.volumeMu.Unlock()
	o.status = Filled
	// o.AppendLog("Order fully filled.")
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// String implements Stringer interface
func (o *Order) String() string {
	return fmt.Sprintf("\norder %s:\n\tside: %s\n\ttype: %s\n\tvolume: %s\n\tprice: %s\n\ttime: %s\n", o.shortOrderID(), o.Side(), o.OrderType(), o.Volume(), o.Price(), o.CreatedAt().String())
}

