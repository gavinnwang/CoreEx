package db

import (
	"database/sql"
	"fmt"
	"github/wry-0313/exchange/internal/config"
	"log"
	// "time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB *sql.DB
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	// log.Println("Waiting 5 seconds for database to start...")
	// time.Sleep(5 * time.Second)
	dsn := buildDSN(cfg)
	conn, err := sql.Open("mysql", dsn)
	log.Printf("Connecting to database...%s\n", dsn)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database ")
	err = conn.Ping() // ping check 
	if err != nil {
		return nil, err
	}
	log.Println("Ping database success")
	// print tables
	rows, err := conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		log.Printf("Table: %s\n", table)
	}
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
