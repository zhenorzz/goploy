package config

import (
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"time"
)

type Config struct {
	Env    string       `toml:"env"`
	APP    APPConfig    `toml:"app"`
	Cookie CookieConfig `toml:"cookie"`
	JWT    JWTConfig    `toml:"jwt"`
	DB     DBConfig     `toml:"db"`
	Log    LogConfig    `toml:"log"`
	Web    WebConfig    `toml:"web"`
	LDAP   LDAPConfig   `toml:"ldap"`
}

type APPConfig struct {
	DeployLimit     int32         `toml:"deployLimit"`
	ShutdownTimeout time.Duration `toml:"shutdownTimeout"`
}

type CookieConfig struct {
	Name   string `toml:"name"`
	Expire int    `toml:"expire"` // second
}

type JWTConfig struct {
	Key string `toml:"key"`
}

type DBConfig struct {
	Type     string `toml:"type"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
}

type LogConfig struct {
	Path  string `toml:"path"`
	Split bool   `toml:"split"`
}

type WebConfig struct {
	Port string `toml:"port"`
}

type LDAPConfig struct {
	Enabled    bool   `toml:"enabled"`
	URL        string `toml:"url"`
	BindDN     string `toml:"bindDN"`
	Password   string `toml:"password"`
	BaseDN     string `toml:"baseDN"`
	UID        string `toml:"uid"`
	UserFilter string `toml:"userFilter"`
}

var Toml Config

func Create(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(config, &Toml)
	if err != nil {
		panic(err)
	}
	setAPPDefault()
	setDBDefault()
}

func setAPPDefault() {
	if Toml.APP.ShutdownTimeout == 0 {
		Toml.APP.ShutdownTimeout = 10
	}
}

func setDBDefault() {
	if Toml.DB.Type == "" {
		Toml.DB.Type = "mysql"
	}
	if Toml.DB.Host == "" {
		Toml.DB.Host = "127.0.0.1"
	}
	if Toml.DB.Port == "" {
		Toml.DB.Port = "3306"
	}
	if Toml.DB.Database == "" {
		Toml.DB.Database = "goploy"
	}
}

func Write(filename string, cfg Config) error {
	yamlData, err := toml.Marshal(&cfg)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}
