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
	Page  int `json:"page"`
	Rows  int `json:"rows"`
	Total int `json:"total"`
}

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

// NewPagination return pagination struct
func NewPagination(param url.Values) (*Pagination, error) {
	page, err := strconv.Atoi(param.Get("page"))
	if err != nil {
		return nil, errors.New("page参数错误")
	}
	rows, err := strconv.Atoi(param.Get("rows"))
	if err != nil {
		return nil, errors.New("rows参数错误")
	}
	pagination := Pagination{Page: page, Rows: rows}
	return &pagination, nil
}
