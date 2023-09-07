package user

import (
	"encoding/json"
	"errors"
	"github/wry-0313/exchange/internal/endpoint"
	"github/wry-0313/exchange/internal/jwt"
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
	userService Service
	jwtService  jwt.Service
	validator   validator.Validate
}

func NewAPI(userService Service, jwtService jwt.Service, validator validator.Validate) *API {
	return &API{
		userService: userService,
		jwtService:  jwtService,
		validator:   validator,
	}
}

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// context := r.Context()

	// Decode request
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("handler: failed to decode request: %v\n", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// Create user and handle errors
	user, err := api.userService.CreateUser(input)
	if err != nil {
		switch {
		case validator.IsValidationError(err):
			endpoint.WriteValidationErr(w, input, err)
		case errors.Is(err, ErrEmailExists):
			endpoint.WriteWithError(w, http.StatusConflict, ErrEmailExists.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}

	jwtToken, err := api.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, CreateUserDTO{User: user, JwtToken: jwtToken})
}

func (api *API) HandleUpdateUserName(w http.ResponseWriter, r *http.Request) {
	// context := r.Context()
	ctx := r.Context()
	userID := middleware.UserIDFromContext(ctx)

	// Decode request
	var input UpdateUserNameInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("handler: failed to decode request: %v\n", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// Update user name and handle errors
	if err := api.userService.UpdateUserName(userID, input.Name); err != nil {
		switch {
		case validator.IsValidationError(err):
			endpoint.WriteValidationErr(w, input, err)
		case errors.Is(err, ErrUserNotFound):
			endpoint.WriteWithError(w, http.StatusNotFound, ErrUserNotFound.Error())
		case errors.Is(err, ErrUserNameSame):
			endpoint.WriteWithError(w, http.StatusConflict, ErrUserNameSame.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, models.SuccessResponse{Message: "User name updated"})
}

func (api *API) HandleGetUserFromJWT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.UserIDFromContext(ctx)
	user, err := api.userService.GetUser(userID)
	if err != nil {
		switch {
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	// write response
	endpoint.WriteWithStatus(w, http.StatusOK, user)
}

// RegisterHandlers is a function that registers all the handlers for the user endpoints
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.HandleCreateUser)
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Post("/name", api.HandleUpdateUserName)
			r.Get("/me", api.HandleGetUserFromJWT)
		})
	})
}
