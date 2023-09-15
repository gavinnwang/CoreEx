package orderbook

import "database/sql"

type repository struct {
	db *sql.DB
}

type Repository interface {
	CreateStock(stock Stock) error
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}