package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"log"
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
	GetUserPrivateInfo(userID string) (UserPrivateInfo, error)

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

	_, err = r.db.Exec("INSERT INTO users (user_id, name, email, password) VALUES (?, ?, ?, ?)", user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	log.Printf("User created successfully: %s\n", user.ID)

	return nil
}

func (r *repository) GetUserPrivateInfo(userID string) (UserPrivateInfo, error) {
	var userPrivateInfo UserPrivateInfo

	rows, err := r.db.Query("SELECT cash_balance FROM users WHERE user_id = ?", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserPrivateInfo{}, ErrUserNotFound
		}
		return UserPrivateInfo{}, fmt.Errorf("repository: failed to get user cash balance: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userPrivateInfo.CashBalance)
		if err != nil {
			return UserPrivateInfo{}, fmt.Errorf("repository: failed to scan user cash balance: %w", err)
		}
	}

	rows, err = r.db.Query("SELECT symbol, volume FROM holdings WHERE user_id = ?", userID)
	if err != nil {
		return UserPrivateInfo{}, fmt.Errorf("repository: failed to get user holdings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var holding models.Holding
		err := rows.Scan(&holding.Symbol, &holding.Volume)
		if err != nil {
			return UserPrivateInfo{}, fmt.Errorf("repository: failed to scan user holdings: %w", err)
		}
		userPrivateInfo.Holdings = append(userPrivateInfo.Holdings, holding)
	}

	rows, err = r.db.Query("SELECT symbol, order_id, order_side, order_status, order_type, filled_at, updated_at, total_processed, volume, initial_volume, price FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return UserPrivateInfo{}, fmt.Errorf("repository: failed to get user orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.Symbol, &order.OrderID, &order.OrderSide, &order.OrderStatus, &order.OrderType, &order.FilledAt, &order.FilledAtTime, &order.TotalProcessed, &order.Volume, &order.InitialVolume, &order.Price)
		if err != nil {
			return UserPrivateInfo{}, fmt.Errorf("repository: failed to scan user orders: %w", err)
		}
		userPrivateInfo.Orders = append(userPrivateInfo.Orders, order)
	}

	return userPrivateInfo, nil
}

// GetUserByEmail returns a single user for a given email.
func (r *repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT user_id, name, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
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
	err := r.db.QueryRow("SELECT user_id, name, email, password FROM users WHERE user_id = ?", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user: %w", err)
	}
	return user, nil
}
