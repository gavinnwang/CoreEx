package orderbook

import (
	"database/sql"
	"fmt"
	"github/wry-0313/exchange/internal/models"
	"log"

	"github.com/shopspring/decimal"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	CreateStock(stock models.Stock) error
	CreateOrder(order *Order, symbol string) error
	UpdateOrder(order *Order, newStatus OrderStatus, newVolume, totalProcessed, filledAt decimal.Decimal) error
	CreateOrUpdateHolding(holding models.Holding) error
	UpdateUserBalance(userID string, newBalance decimal.Decimal) error
	CreateMarketPriceHistory(symbol string, priceHistory models.StockPriceHistory) error
	GetEntireMarketPriceHistory(symbol string) ([]models.StockPriceHistory, error)
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

	orderSide := order.side.String()
	orderStatus := order.status.String()

	sql := `INSERT INTO orders (user_id, order_id, order_side, order_status, order_type, volume, initial_volume, price, created_at, symbol) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(sql, order.userID.String(), order.orderID.String(), orderSide, orderStatus, order.orderType.String(), order.volume, order.volume, order.price, order.createdAt, symbol)
	if err != nil {
		return err
	}

	log.Printf("Order created successfully: %s\n", order.shortOrderID())

	return nil
}

func (r *repository) UpdateOrder(order *Order, newStatus OrderStatus, newVolume, totalProcessed, filledAt decimal.Decimal) error {

	sql := `UPDATE orders SET order_status = ?, volume = ?, filled_at = ?, total_processed = total_processed + ? WHERE order_id = ?`

	_, err := r.db.Exec(sql, newStatus.String(), newVolume, filledAt, totalProcessed, order.orderID.String())
	if err != nil {
		return fmt.Errorf("repository: failed to update order: %v", err)
	}

	log.Printf("Order updated successfully: %s\n", order.shortOrderID())

	return nil
}

func (r *repository) CreateOrUpdateHolding(holding models.Holding) error {

	log.Printf("holding: %+v\n", holding)

	sql := `CALL InsertOrUpdateHoldingThenDeleteZeroVolume(?, ?, ?)`

	_, err := r.db.Exec(sql, holding.UserID, holding.Symbol, holding.VolumeChange)
	if err != nil {
		return fmt.Errorf("repository: failed to create or update holding: %v", err)
	}

	log.Printf("Holding created or updated successfully: %s\n", holding.UserID[22:])

	return nil
}

func (r *repository) UpdateUserBalance(userID string, balanceChange decimal.Decimal) error {

	sql := `UPDATE users SET cash_balance = cash_balance + ? WHERE user_id = ?`

	_, err := r.db.Exec(sql, balanceChange, userID)
	if err != nil {
		return fmt.Errorf("repository: failed to update user balance: %v", err)
	}

	log.Printf("User balance updated successfully: %s\n", userID[22:])

	return nil
}

func (r *repository) CreateMarketPriceHistory(symbol string, priceHistory models.StockPriceHistory) error {

	sql := `INSERT INTO stock_history (symbol, open, high, low, close, volume, recorded_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(sql, symbol, priceHistory.Open, priceHistory.High, priceHistory.Low, priceHistory.Close, priceHistory.Volume, priceHistory.RecordedAt)

	if err != nil {
		return fmt.Errorf("repository: failed to create market price history: %v", err)
	}

	return nil
}

func (r *repository) GetEntireMarketPriceHistory(symbol string) ([]models.StockPriceHistory, error) {

	sql := `SELECT open, high, low, close, volume, recorded_at 
	FROM (
		SELECT open, high, low, close, volume, recorded_at 
		FROM stock_history 
		WHERE symbol = ? 
		ORDER BY recorded_at DESC 
		LIMIT 75
	) AS sub
	ORDER BY recorded_at ASC;
	`

	rows, err := r.db.Query(sql, symbol)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get market price history: %v", err)
	}
	defer rows.Close()

	var priceData []models.StockPriceHistory
	for rows.Next() {
		var history models.StockPriceHistory
		if err := rows.Scan(&history.Open, &history.High, &history.Low, &history.Close, &history.Volume, &history.RecordedAt); err != nil {
			return nil, fmt.Errorf("repository: failed to scan row: %v", err)
		}
		priceData = append(priceData, history)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error iterating rows: %v", err)
	}

	return priceData, nil
}
