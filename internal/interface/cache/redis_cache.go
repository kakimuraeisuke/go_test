package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go_test/internal/usecase"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache はCacheインターフェースを実装します
type redisCache struct {
	client *redis.Client
}

// NewRedisCache は新しいRedisキャッシュを作成します
func NewRedisCache(client *redis.Client) usecase.Cache {
	return &redisCache{client: client}
}

// Set はRedisキャッシュに値を設定します
func (c *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = c.client.Set(ctx, key, jsonValue, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get はRedisキャッシュから値を取得します
func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key not found")
		}
		return "", fmt.Errorf("failed to get cache: %w", err)
	}

	return val, nil
}

// Delete はRedisキャッシュから値を削除します
func (c *redisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}
