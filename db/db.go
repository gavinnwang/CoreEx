package db

import "database/sql"

type db struct {
	db *sql.DB
}

// func (db *db) New() *db {
// 	return &db{
// 		db:
// 	}
// }