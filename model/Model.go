package model

import (
	"database/sql"
	"log"
	"os"
)

// NewDB return db model
func NewDB() *sql.DB {
	dbType := os.Getenv("DB_TYPE")
	dbConn := os.Getenv("DB_CONN")
	db, err := sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
