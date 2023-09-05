package ws

import (
	"encoding/json"
	"github/wry-0313/exchange/exchange"
)

const (

	// Events
	EventStreamMarketPrice = "exchange.stream_price"

	EventStreamSymbolInfo = "exchange.stream_info"

	// CloseReasonBadEvent indicates that the event field has an incorrect type.
	CloseReasonBadEvent = "The event field is an incorrect type."

	// CloseReasonInternalServer indicates an internal server error.
	CloseReasonInternalServer = "Internal server error."

	// CloseReasonUnsupportedEvent indicates that the event is unsupported.
	CloseReasonUnsupportedEvent = "The event is unsupported."

	// CloseReasonBadParams indicates that the params have incorrect field types.
	CloseReasonBadParams = "The params have incorrect field types."

	// CloseReasonUnauthorized indicates an unauthorized request.
	CloseReasonUnauthorized = "Unauthorized."

	// ErrMsgInternalServer indicates an internal server error.
	ErrMsgInternalServer = "Internal server error."
)

// Request is a struct that describes the shape of every message request.
type Request struct {
	Event  string          `json:"event"`
	Params json.RawMessage `json:"params"`
}

// ResponseBase represents the base response structure.
type ResponseBase struct {
	Event        string `json:"event"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ParamsSymbol contains the parameter symbol and is used for handlers that only needs the symbol
type ParamsSymbol struct {
	Symbol string `json:"symbol" validate:"required"`
}

type ResponseGetMarketPrice struct {
	ResponseBase
	Result ResultGetMarketPrice `json:"result,omitempty"`
}

type ResponseGetSymbolInfo struct {
	ResponseBase
	Result exchange.SymbolInfoResponse `json:"result,omitempty"`
}

// ResultBoardConnect contains the result of board connection.
type ResultGetMarketPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
