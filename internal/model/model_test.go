package model

import (
	"database/sql"
	"fmt"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/database"
	"testing"
)

func TestInitDB(t *testing.T) {
	config.Init()
	db, err := sql.Open(config.Toml.DB.Type, fmt.Sprintf(
		"%s:%s@(%s:%s)/?charset=utf8mb4,utf8\n",
		config.Toml.DB.User,
		config.Toml.DB.Password,
		config.Toml.DB.Host,
		config.Toml.DB.Port))
	if err != nil {
		t.Fatal(err)
	}
	if err := CreateDB(db, config.Toml.DB.Database); err != nil {
		t.Fatal(err)
	}
	if err := UseDB(db, config.Toml.DB.Database); err != nil {
		t.Fatal(err)
	}
	if err := ImportSQL(db, database.GoploySQL); err != nil {
		t.Fatal(err)
	}
}
