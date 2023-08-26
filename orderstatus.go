package orderbook

type OrderStatus int

const (
	Rejected OrderStatus = iota
	Received
	Open
	Cancelled
	Filled
)

// String implements fmt.Stringer interface
func (os OrderStatus) String() string {
	switch os {
	case Rejected:
		return "Rejected"
	case Received:
		return "Received"
	case Open:
		return "Open"
	case Cancelled:
		return "Cancelled"
	case Filled:
		return "Filled"
	default:
		return "Unknown"
	}
}
