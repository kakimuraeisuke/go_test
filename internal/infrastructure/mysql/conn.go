package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds MySQL configuration
type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// NewConfig creates a new MySQL config from environment variables
func NewConfig() *Config {
	return &Config{
		Host:     getEnv("MYSQL_HOST", "localhost"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		Database: getEnv("MYSQL_DATABASE", "go_test"),
		User:     getEnv("MYSQL_USER", "user"),
		Password: getEnv("MYSQL_PASSWORD", "hogehoge"),
	}
}

// Connect creates a new MySQL connection
func Connect(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
