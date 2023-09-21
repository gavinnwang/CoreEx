package models

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID        ulid.ULID `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email"`
	Password  *string   `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
