package auth

import (
	"encoding/json"
	"errors"
	"github/wry-0313/exchange/endpoint"
	"github/wry-0313/exchange/pkg/validator"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	// ErrMsgInternalServer is a message displayed when an unexpected error occurs
	ErrMsgInternalServer = "Internal server error"
)

type API struct {
	authService Service 
	validator validator.Validate
}

// NewAPI creates a new intance of the API struct.
func NewAPI(authService Service, validator validator.Validate) API {
	return API{
		authService: authService,
		validator:   validator,
	}
}

// HandleLogin handles a user's login request. It returns a token in the response
// if the login is successful.
func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()
	token, err := api.authService.Login(input)
	if err != nil {
		switch {
		case errors.Is(err, errBadLogin):
			endpoint.WriteWithError(w, http.StatusUnauthorized, errBadLogin.Error())
		case validator.IsValidationError(err):
			endpoint.WriteValidationErr(w, input, err)
		default:
			log.Printf("HandleLogin: Failed to login user due to internal server error: %v", err)
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, LoginDTO{Token: token})
}

// RegisterHandlers registers the API's request handlers.
func (api *API) RegisterHandlers(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", api.HandleLogin)
	})
}