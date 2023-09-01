package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize database connection
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS test_db")
	if err != nil {
		log.Fatal(err)
	}

	// Use the created database
	_, err = db.Exec("USE test_db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(64) NOT NULL,
		password VARCHAR(64) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ('john_doe', 'password123')")
	if err != nil {
		log.Fatal(err)
	}

	// Query data
	rows, err := db.Query("SELECT id, username, password FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through query results
	for rows.Next() {
		var id int
		var username, password string

		if err := rows.Scan(&id, &username, &password); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %d, Username: %s, Password: %s\n", id, username, password)
	}
}
