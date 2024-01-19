package config

type CookieConfig struct {
	Name   string `toml:"name"`
	Expire int    `toml:"expire"` // second
}
