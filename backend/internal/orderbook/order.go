package orderbook

import (
	// "encoding/json"
	// "log"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"log"
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
	// totalProcessed decimal.Decimal
	createdAt time.Time
	volumeMu  sync.RWMutex
}

func (s *service) NewOrder(side Side, userID ulid.ULID, orderType OrderType, price, volume decimal.Decimal, partialAllowed bool) *Order {
	loc, _ := time.LoadLocation("America/Chicago")
	o := &Order{
		side:      side,
		orderID:   ulid.Make(),
		userID:    userID,
		orderType: orderType,
		status:    Open,
		price:     price,
		volume:    volume,
		createdAt: time.Now().In(loc),
	}
	go func() {
		err := s.obRepo.CreateOrder(o, s.symbol)
		if err != nil {
			log.Fatalf("service: failed to create order: %v", err)
		}
	}()
	return o
}

// ID returns orderID field copy
func (o *Order) OrderID() ulid.ULID {
	return o.orderID
}

// shortOrderID returns last 4 characters of orderID (for debugging purposes)
func (o *Order) shortOrderID() string {
	return o.orderID.String()[22:]
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

func (s *service) fillOrder(o *Order, filledVolume, filledAt decimal.Decimal) {
	// log.Printf("service: order %s filled with volume %s at price %s\n", o.shortOrderID(), filledVolume, filledAt)
	o.volumeMu.Lock()
	newVolume := o.volume.Sub(filledVolume)
	o.volume = newVolume
	o.volumeMu.Unlock()
	if newVolume.IsZero() {
		o.status = Filled
	} else {
		o.status = PartiallyFilled
	}
	go func() {
		processedValue := filledVolume.Mul(filledAt).Round(2)
		err := s.obRepo.UpdateOrder(o, o.status, newVolume, processedValue, filledAt)
		if err != nil {
			log.Fatalf("service: failed to update order: %v", err)
		}
		var holdingChange models.HoldingChange
		if o.Side() == Buy { // order wants to buy stock so we need to add new holding to user and subtract user balance
			holdingChange = models.HoldingChange{
				UserID:       o.UserID().String(),
				Symbol:       s.symbol,
				VolumeChange: filledVolume,
			}
			processedValue = processedValue.Neg()
		} else {
			holdingChange = models.HoldingChange{
				UserID:       o.UserID().String(),
				Symbol:       s.symbol,
				VolumeChange: filledVolume.Neg(),
			}
		}
		err = s.obRepo.CreateOrUpdateHolding(holdingChange)
		if err != nil {
			log.Fatalf("service: failed to create or update holding: %v", err)
		}
		err = s.obRepo.UpdateUserBalance(o.UserID().String(), processedValue)
		if err != nil {
			log.Fatalf("service: failed to update user balance: %v", err)
		}
	}()	
}


func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// String implements Stringer interface
func (o *Order) String() string {
	return fmt.Sprintf("\norder %s:\n\tside: %s\n\ttype: %s\n\tvolume: %s\n\tprice: %s\n\ttime: %s\n", o.shortOrderID(), o.Side(), o.OrderType(), o.Volume(), o.Price(), o.CreatedAt().String())
}
