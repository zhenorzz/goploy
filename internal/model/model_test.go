package model

import (
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/database"
	"testing"
)

func TestInitDB(t *testing.T) {
	config.Init()
	db, err := Open(config.Toml.DB)

	if err != nil {
		t.Fatal(err)
	}
	if err := db.CreateDB(config.Toml.DB.Database); err != nil {
		t.Fatal(err)
	}
	if err := db.UseDB(config.Toml.DB.Database); err != nil {
		t.Fatal(err)
	}
	if err := db.ImportSQL(database.GoploySQL); err != nil {
		t.Fatal(err)
	}
}
