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

	_, err = r.db.Exec("INSERT INTO stocks (symbol, name) VALUES (?, ?)", stock.Symbol, stock.Name)
	if err != nil {
		return err
	}

	log.Printf("Stock created successfully: %s\n", stock.Symbol)
	return nil
}
