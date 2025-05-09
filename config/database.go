package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig holds all database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type connectionString struct {
	conn DatabaseConfig `toml:"database"`
}

func LoadDataFromToml(filepath string) (*DatabaseConfig, error) {
	var config connectionString
	_, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		log.Println("Error loading config file:", err)
		return nil, err
	}
	return &config.conn, nil
}

// GetDBConfig returns database configuration
func GetDBConfig() *DatabaseConfig {
	conn, err := LoadDataFromToml("/test.toml")
	if conn == nil {
		log.Println("Error loading config:", err)
	}
	log.Println("Database Connection", conn)
	return conn
}

// GetDSN returns the Data Source Name
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func InitDB() {
	log.Println("Initializing database connection...")
	dbConfig := GetDBConfig()
	dsn := dbConfig.GetDSN()
	log.Printf("Connecting to database with DSN: %s", dsn)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connection established successfully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
