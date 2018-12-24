package model

import (
	"database/sql"
	"log"
	"os"
	"fmt"
)

func NewDB() *sql.DB {
	dbType := os.Getenv("DB_TYPE")
	dbConn := os.Getenv("DB_CONN")
	fmt.Println(dbConn,dbType)

	db, err := sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
