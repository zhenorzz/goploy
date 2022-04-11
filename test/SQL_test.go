package test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/model"
	"testing"
)

func TestInitDB(t *testing.T) {
	config.Create("../goploy.toml")
	db, err := sql.Open(config.Toml.DB.Type, fmt.Sprintf(
		"%s:%s@(%s:%s)/?charset=utf8mb4,utf8\n",
		config.Toml.DB.User,
		config.Toml.DB.Password,
		config.Toml.DB.Host,
		config.Toml.DB.Port))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := model.CreateDB(db, config.Toml.DB.Database); err != nil {
		t.Fatal(err)
	}
	if err := model.UserDB(db, config.Toml.DB.Database); err != nil {
		t.Fatal(err)
	}
	if err := model.ImportSQL(db, "sql/goploy.sql"); err != nil {
		t.Fatal(err)
	}
}
