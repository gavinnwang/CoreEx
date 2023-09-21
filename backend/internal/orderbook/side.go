package orderbook

// import (
// 	"encoding/json"
// 	"reflect"
// )

type Side int

const (
	Sell Side = iota
	Buy
	Invalid
)

// String implements fmt.Stringer interface
func (s Side) String() string {
	if s == Buy {
		return "Buy"
	}
	return "Sell"
}

func SideFromString(s string) (Side, error) {

	switch s {
	case "buy":
		return Buy, nil
	case "sell":
		return Sell, nil
	default:
		return Invalid, ErrInvalidSide
	}
}

// // MarshalJSON implements json.Marshaler interface
// func (s Side) MarshalJSON() ([]byte, error) {
// 	return []byte(`"` + s.String() + `"`), nil
// }

// // UnmarshalJSON implements json.Unmarshaler interface
// func (s *Side) UnmarshalJSON(data []byte) error {
// 	switch string(data) {
// 	case `"buy"`:
// 		*s = Buy
// 	case `"sell"`:
// 		*s = Sell
// 	default:
// 		return &json.UnsupportedValueError{
// 			Value: reflect.New(reflect.TypeOf(data)),
// 			Str:   string(data),
// 		}
// 	}
// 	return nil
// }
