package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github/wry-0313/exchange/models"
)

var (
	ErrEmailExists          = errors.New("User with this email already exists")
	ErrUserNotFound         = errors.New("User does not exist")
	ErrVerificationNotFound = errors.New("Verification does not exist")
	ErrUserNameSame         = errors.New("User name is the same")
)

type Repository interface {
	CreateUser(user models.User) error

	GetUser(userID string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	UpdateUserName(userID, name string) error
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
		return fmt.Errorf("repository: failed to check if user exists: %w", err)
	}
	if exists {
		return ErrEmailExists // Email already exists
	}

	_, err = r.db.Exec("INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail returns a single user for a given email.
func (r *repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user: %w", err)
	}
	return user, nil
}

func (r *repository) UpdateUserName(userID, name string) error {
	user, err := r.GetUser(userID)
	if err != nil {
		return err
	}
	if user.Name == name {
		return ErrUserNameSame
	}

	_, err = r.db.Exec("UPDATE users SET name = ? WHERE id = ?", name, userID)
	if err != nil {
		return fmt.Errorf("repository: failed to update user name: %w", err)
	}
	return nil
}

func (r *repository) GetUser(userID string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user: %w", err)
	}
	return user, nil
}
