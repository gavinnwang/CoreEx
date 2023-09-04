package db

import (
	"database/sql"
	"fmt"
	"github/wry-0313/exchange/config"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB *sql.DB
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	dsn := buildDSN(cfg)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connected to database %s\n", dsn)
	return &DB{
		DB: conn,
	}, nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%v?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}