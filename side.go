package orderbook

import (
	"encoding/json"
	"reflect"
)

type Side int

const (
	Ask Side = iota
	Bid
)

// String implements fmt.Stringer interface
func (s Side) String() string {
	if s == Bid {
		return "bid"
	}
	return "ask"
}

// MarshalJSON implements json.Marshaler interface
func (s Side) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (s *Side) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"buy"`:
		*s = Bid
	case `"sell"`:
		*s = Ask
	default:
		return &json.UnsupportedValueError{
			Value: reflect.New(reflect.TypeOf(data)),
			Str:   string(data),
		}
	}
	return nil
}
