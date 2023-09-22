package orderbook

import (
	"database/sql"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"log"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	CreateStock(stock models.Stock) error
	CreateOrder(order *Order, symbol string) error
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateStock(stock models.Stock) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM stocks WHERE symbol = ?)", stock.Symbol).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("repository: failed to check if stock exists: %w", err)
	}
	if exists {
		log.Printf("Stock already exists: %s\n", stock.Symbol)
		return nil
	}

	_, err = r.db.Exec("INSERT INTO stocks (symbol) VALUES (?)", stock.Symbol)
	if err != nil {
		return err
	}

	log.Printf("Stock created successfully: %s\n", stock.Symbol)
	return nil
}

func (r *repository) CreateOrder(order *Order, symbol string) error {
	fmt.Printf("Order: %v\n", order.OrderID())
	orderSide := order.side.String()
	orderStatus := order.status.String()

	sql := `INSERT INTO orders (user_id, order_id, order_side, order_status, volume, price, created_at, symbol) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	log.Printf("User id %s\n", order.userID.String())
	_, err := r.db.Exec(sql, order.userID.String(), order.orderID.String(), orderSide, orderStatus, order.volume, order.price, order.createdAt, symbol)
	if err != nil {
		return err
	}

	log.Printf("Order created successfully: %s\n", order.orderID)

	return nil
}
