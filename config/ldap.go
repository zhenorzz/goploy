package config

type LDAPConfig struct {
	Enabled    bool   `toml:"enabled"`
	URL        string `toml:"url"`
	BindDN     string `toml:"bindDN"`
	Password   string `toml:"password"`
	BaseDN     string `toml:"baseDN"`
	UID        string `toml:"uid"`
	Name       string `toml:"name"`
	UserFilter string `toml:"userFilter"`
}
