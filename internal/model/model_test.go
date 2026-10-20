package model

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/database"
)

// configFilePath is the same path config.Init loads, without panicking when
// the test binary cannot be resolved.
func configFilePath(t *testing.T) string {
	t.Helper()
	dir := config.AssetDir
	if dir == "" {
		file, err := exec.LookPath(os.Args[0])
		if err != nil {
			t.Skipf("goploy.toml location is unknown: LookPath(%s): %v", os.Args[0], err)
		}
		abs, err := filepath.Abs(file)
		if err != nil {
			t.Skipf("goploy.toml location is unknown: %v", err)
		}
		dir = filepath.Dir(abs)
	}
	return filepath.Join(dir, "goploy.toml")
}

func TestInitDB(t *testing.T) {
	cfg := configFilePath(t)
	if _, err := os.Stat(cfg); err != nil {
		t.Skipf("goploy.toml is missing at %s: %v", cfg, err)
	}

	config.Init()

	db, err := Open(config.Toml.DB)
	if err != nil {
		t.Skipf("database is not reachable: %v", err)
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
