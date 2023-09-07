package exchange

import (
	"encoding/json"
	"errors"
	"github/wry-0313/exchange/internal/endpoint"
	"github/wry-0313/exchange/internal/middleware"
	"github/wry-0313/exchange/internal/models"
	"github/wry-0313/exchange/pkg/validator"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	ErrMsgInternalServer = "Internal server error"
)

type API struct {
	exchangeService Service
}

func NewAPI(exchangeService Service) *API {
	return &API{
		exchangeService: exchangeService,
	}
}

func (api *API) HandlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.UserIDFromContext(ctx)
	log.Printf("API: user requests to place order: %v\n", userID[:4])

	var input PlaceOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("handler: failed to decode request: %v\n", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}

	// Set userID
	input.UserID = userID

	defer r.Body.Close()

	if err := api.exchangeService.PlaceOrder(input); err != nil {
		switch {
		case validator.IsValidationError(err):
			endpoint.WriteValidationErr(w, input, err)
		case errors.Is(err, ErrInvalidSymbol):
			endpoint.WriteWithError(w, http.StatusBadRequest, err.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, models.SuccessResponse{Message: "Order placed"})
}

// func (api *API) HandleStreamMarketPrice(w http.ResponseWriter, r *http.Request) {
// 	// Unmarshal params
// 	var params StreamPriceParams
// 	if err := ws.UnmarshalParams(msgReq, &params, c); err != nil {
// 		return
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
// 		log.Printf("handler: failed to decode request: %v\n", err)
// 		endpoint.HandleDecodeErr(w, err)
// 		return
// 	}

// 	conn, err := ws.Upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			p, err := api.exchangeService.GetMarketPrice(input.Symbol)
// 			if err != nil {
// 				log.Printf("handler: failed to get market price: %v\n", err)
// 				endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
// 				return
// 			}
// 			priceString := p.String()

// 			if err := conn.WriteMessage(websocket.TextMessage, []byte(priceString)); err != nil {
// 				log.Println("WriteMessage Error:", err)
// 				return
// 			}
// 		}
// 	}
// }

// RegisterHandlers is a function that registers all the handlers for the user endpoints
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {

	r.Route("/orders", func(r chi.Router) {
		// r.HandleFunc("/price", api.HandleStreamMarketPrice)
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Post("/", api.HandlePlaceOrder)
		})
	})
}
