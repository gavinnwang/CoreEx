package orderbook

type SymbolInfoResponse struct {
	Symbol    string  `json:"symbol"`
	AskVolume float64 `json:"ask_volume"`
	BidVolume float64 `json:"bid_volume"`
	BestBid   float64 `json:"best_bid"`
	BestAsk   float64 `json:"best_ask"`
	Price     float64 `json:"price"`
}