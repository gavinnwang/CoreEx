package exchange

import (
	"encoding/json"
	"fmt"
	"github/wry-0313/exchange/endpoint"
	"github/wry-0313/exchange/middleware"
	"github/wry-0313/exchange/models"
	"github/wry-0313/exchange/pkg/validator"
	ws "github/wry-0313/exchange/websocket"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
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
	fmt.Printf("user requests to place order: %v\n", userID[:4])

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
		case err == ErrInvalidSymbol:
			endpoint.WriteWithError(w, http.StatusNoContent, err.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, models.MessageResponse{Message: "Order placed"})
}

func (api *API) HandleStreamMarketPrice(w http.ResponseWriter, r *http.Request) {
	var input StreamPriceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("handler: failed to decode request: %v\n", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}

	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p, err := api.exchangeService.GetMarketPrice(input.Symbol)
			if err != nil {
				log.Printf("handler: failed to get market price: %v\n", err)
				endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
				return
			}
			priceString := p.String()

			if err := conn.WriteMessage(websocket.TextMessage, []byte(priceString)); err != nil {
				log.Println("WriteMessage Error:", err)
				return
			}
		}
	}
}

// RegisterHandlers is a function that registers all the handlers for the user endpoints
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {

	r.Route("/order", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Post("/", api.HandlePlaceOrder)
		})
	})
}
