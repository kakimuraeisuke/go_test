package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Config はRedis設定を保持します
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewConfig は環境変数から新しいRedis設定を作成します
func NewConfig() *Config {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	return &Config{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
}

// Connect は新しいRedis接続を作成します
func Connect(config *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return rdb
}

// Ping はRedis接続をテストします
func Ping(ctx context.Context, client *redis.Client) error {
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}
	return nil
}

// getEnv はデフォルト値付きで環境変数を取得します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
