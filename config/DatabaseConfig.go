package config

type DatabaseConnection struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"` // or int if you used an integer in toml
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	SSLMode  string `toml:"sslmode"`
}
