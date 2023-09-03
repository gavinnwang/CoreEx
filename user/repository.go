package user

import (
	"database/sql"
	"errors"
	"fmt"
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
	GetUserByEmail(email string) (models.User, error)

	// DeleteUser(userID uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
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
	fmt.Printf("User: %v\n", user)

	_, err = r.db.Exec("INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail returns a single user for a given email.
func (r *repository) GetUserByEmail( email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}