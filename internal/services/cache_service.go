package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ssjoy/website-health-checker/internal/redis"
)

// CacheService defines the operations for key-value store interactions.
type CacheService interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type cacheService struct {
	rdb *redis.Client
}

// NewCacheService creates a new instance of CacheService.
func NewCacheService(rdb *redis.Client) CacheService {
	return &cacheService{rdb: rdb}
}

// Set saves a key-value pair into Redis cache with optional expiration.
func (s *cacheService) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	err := s.rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Get retrieves a key's value from Redis cache.
func (s *cacheService) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}
	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
