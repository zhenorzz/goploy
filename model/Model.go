package model

import (
	"database/sql"
	"log"
	"os"
)

func NewDB() *sql.DB {
	dbType := os.Getenv("DB_TYPE")
	dbConn := os.Getenv("DB_CONN")
	db, err := sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
