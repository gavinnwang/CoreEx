package orderbook

type OrderType int

const (
	Limit = iota
	Market
	Stop
)

// String implements fmt.Stringer interface
func (ot OrderType) String() string {
	if ot == Limit {
		return "Limit"
	}
	return "Market"
}
