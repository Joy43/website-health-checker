package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ssjoy/website-health-checker/internal/config"
)

// Client wraps the redis.Client to expose health check and pool control methods.
type Client struct {
	*redis.Client
}

// Connect creates a new Redis client based on application configuration.
func Connect(cfg *config.Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		PoolSize:     20,
		MinIdleConns: 10,
	})

	// Perform initial ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Client{rdb}, nil
}

// CheckHealth queries Redis ping status.
func (c *Client) CheckHealth(ctx context.Context) error {
	return c.Ping(ctx).Err()
}
