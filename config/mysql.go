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
	d.SetDefault()
	return nil
}

func (d *DBConfig) SetDefault() {
	if d.Type == "" {
		d.Type = "mysql"
	}
	if d.Host == "" {
		d.Host = "127.0.0.1"
	}
	if d.Port == "" {
		d.Port = "3306"
	}
	if d.Database == "" {
		d.Database = "goploy"
	}
}
