package config

type DBConfig struct {
	Type     string `toml:"type"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
}

func (d *DBConfig) OnChange() error {
	setDBDefault()
	return nil
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
