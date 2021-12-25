package model

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/utils"
	"log"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
)

//go:embed sql
var sqlFile embed.FS

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

// review state type
const (
	PENDING = iota
	APPROVE
	DENY
)

// DB init when the program start
var DB *sql.DB

// Init DB
func Init() {
	dbConn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4,utf8",
		config.Toml.DB.User,
		config.Toml.DB.Password,
		config.Toml.DB.Host,
		config.Toml.DB.Port,
		config.Toml.DB.Database,
	)
	var err error
	DB, err = sql.Open(config.Toml.DB.Type, dbConn)
	if err != nil {
		log.Fatal(err)
	}
}

// PaginationFrom param return pagination struct
func PaginationFrom(param url.Values) (Pagination, error) {
	page, err := strconv.ParseUint(param.Get("page"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("invalid page")
	}
	rows, err := strconv.ParseUint(param.Get("rows"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("invalid rows")
	}
	pagination := Pagination{Page: page, Rows: rows}
	return pagination, nil
}

// ImportSQL -
func ImportSQL(db *sql.DB, sqlPath string) error {
	sqlContent, err := sqlFile.ReadFile(sqlPath)
	if err != nil {
		return err
	}
	for _, query := range strings.Split(string(sqlContent), ";") {
		query = utils.ClearNewline(query)
		if len(query) == 0 {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func Update(targetVerStr string) error {
	systemConfig, err := SystemConfig{
		Key: "version",
	}.GetDataByKey()
	if err != nil {
		return err
	}

	if systemConfig.Value == "" {
		systemConfig.Value = "0.0.1"
	}

	currentVer, err := version.NewVersion(systemConfig.Value)
	if err != nil {
		return err
	}

	targetVer, err := version.NewVersion(targetVerStr)
	if err != nil {
		return err
	}

	if ret := currentVer.Compare(targetVer); ret == 0 {
		return nil
	} else if ret == 1 {
		return errors.New("currentVer greater than targetVer")
	}

	sqlEntries, err := sqlFile.ReadDir("sql")
	if err != nil {
		return err
	}
	var vers []*version.Version
	for _, entry := range sqlEntries {
		filename := entry.Name()
		ver, err := version.NewVersion(filename[0 : len(filename)-len(path.Ext(filename))])
		if err != nil {
			continue
		}
		vers = append(vers, ver)
	}

	sort.Sort(version.Collection(vers))

	for _, ver := range vers {
		if currentVer.LessThan(ver) && targetVer.GreaterThanOrEqual(ver) {
			if err := ImportSQL(DB, "sql/"+ver.String()+".sql"); err != nil {
				return err
			}
		}
	}
	println(`Update app success`)
	systemConfig.Value = targetVerStr
	return systemConfig.EditRowByKey()
}
