package config

import (
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
)

type Config struct {
	Env    string       `toml:"env"`
	Cookie CookieConfig `toml:"cookie"`
	JWT    JWTConfig    `toml:"jwt"`
	DB     DBConfig     `toml:"db"`
	Log    LogConfig    `toml:"log"`
	Web    WebConfig    `toml:"web"`
	LDAP   LDAPConfig   `toml:"ldap"`
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
}

type LogConfig struct {
	Path string `toml:"path"`
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
