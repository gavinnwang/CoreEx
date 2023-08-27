package orderbook

type OrderStatus int

const (
	Rejected OrderStatus = iota
	Open
	Cancelled
	Filled
	PartiallyFilled
)

// String implements fmt.Stringer interface
func (os OrderStatus) String() string {
	switch os {
	case Rejected:
		return "Rejected"
	case Open:
		return "Open"
	case Cancelled:
		return "Cancelled"
	case Filled:
		return "Filled"
	case PartiallyFilled:
		return "PartiallyFilled"
	default:
		return "Unknown"
	}
}
