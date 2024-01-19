package config

type CORSConfig struct {
	Enabled     bool   `toml:"enabled"`
	Origins     string `toml:"origins"`
	Methods     string `toml:"methods"`
	Headers     string `toml:"headers"`
	Credentials bool   `toml:"credentials"`
}
