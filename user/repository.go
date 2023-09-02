package user

import (
	"database/sql"
	"errors"
	"github/wry-0313/exchange/models"
	// "github.com/google/uuid"
)

var (
	ErrEmailExists = errors.New("User with this email already exists")
	ErrUserNotFound = errors.New("User does not exist")
	ErrVerificationNotFound = errors.New("Verification does not exist")
)

type Repository interface {
	CreateUser(user models.User) error

	// GetUser(userID uuid.UUID) (models.User, error)
	// GetUserByEmail(email string) (models.User, error)

	// DeleteUser(userID uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(user models.User) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err // something went wrong with the query
	}
	if exists {
		return ErrEmailExists // Email already exists
	}

	_, err = r.db.Exec("INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
