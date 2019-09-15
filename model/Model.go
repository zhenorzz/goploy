package model

import (
	"database/sql"
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"
)

// Pagination struct
type Pagination struct {
	Page  uint64 `json:"page"`
	Rows  uint64 `json:"rows"`
	Total uint64 `json:"total"`
}

// state type
const (
	Fail = iota
	Success
)

// state type
const (
	Disable = iota
	Enable
)

// DB init when the program start
var DB *sql.DB

// Init DB open
func Init() {
	dbType := os.Getenv("DB_TYPE")
	dbConn := os.Getenv("DB_CONN")
	var err error
	DB, err = sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatal(err)
	}
}

// PaginationFrom param return pagination struct
func PaginationFrom(param url.Values) (Pagination, error) {
	page, err := strconv.ParseUint(param.Get("page"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("Invalid page")
	}
	rows, err := strconv.ParseUint(param.Get("rows"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("Invalid rows")
	}
	pagination := Pagination{Page: page, Rows: rows}
	return pagination, nil
}
