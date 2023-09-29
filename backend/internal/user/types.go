package user

import "github/wry-0313/exchange/internal/models"

// CreateUserInput defines the structure for requests to create a new user.
type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=24"`
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"password" validate:"min=8"`
}

type UpdateUserNameInput struct {
	Name string `json:"name" validate:"required,min=2,max=24"`
}

// CreateUserDTO defines the structure of a successful create user response.
type CreateUserDTO struct {
	User     models.User `json:"user"`
	JwtToken string      `json:"jwt_token"`
}

type UserPrivateInfo struct {
	CashBalance float64          `json:"cash_balance"`
	Holdings    []models.Holding `json:"holdings"`
	Orders      []models.Order   `json:"orders"`
}
