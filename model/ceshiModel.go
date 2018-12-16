package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func NewDB() *sql.DB {
	godotenv.Load()
	dbType := os.Getenv("DB_TYPE")
	dbConn := os.Getenv("DB_CONN")
	db, err := sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Query() {
	db := NewDB()
	rows, err := db.Query("SELECT name FROM ceshi")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
