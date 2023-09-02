package user

// CreateUserInput defines the structure for requests to create a new user.
type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=24"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
}