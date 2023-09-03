package user

import "github/wry-0313/exchange/models"

// CreateUserInput defines the structure for requests to create a new user.
type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=24"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
}

// CreateUserDTO defines the structure of a successful create user response.
type CreateUserDTO struct {
	User     models.User `json:"user"`
	JwtToken string      `json:"jwt_token"`
}