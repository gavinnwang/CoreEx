package orderbook

import "github/wry-0313/exchange/internal/models"

type SymbolInfoResponse struct {
	Symbol     string     `json:"symbol"`
	AskVolume  float64    `json:"ask_volume"`
	BidVolume  float64    `json:"bid_volume"`
	BestBid    float64    `json:"best_bid"`
	BestAsk    float64    `json:"best_ask"`
	Price      float64    `json:"price"`
	CandleData CandleData `json:"candle_data"`
}

type CandleData struct {
	models.StockPriceHistory
	NewCandle bool `json:"new_candle"`
}

type RedisPubMsgBase struct {
	Event        string `json:"event"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type SymbolInfoPubMsg struct {
	RedisPubMsgBase
	Result SymbolInfoResponse `json:"result,omitempty"`
}

const (
	EventStreamSymbolInfo = "exchange.stream_info"
)
