package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Operations defines the interface for interacting with Redis
type Operations interface {
	Set(ctx context.Context, key string, value interface{}) error
	SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetList(ctx context.Context, key string, values []string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}

// Redis wraps a Redis client
type Redis struct {
	*redis.Client
}

// NewRedis initializes a new Redis client based on the given configuration
//
// Example:
//
//	cfg := &Config{}
//	redisClient, err := NewRedis(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewRedis(cfg *Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Redis")
	return &Redis{client}, nil
}

// Set stores a key-value pair without expiration
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.Set(ctx, "user:123", "John Doe")
func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}

// SetWithExpiration stores a key-value pair with expiration
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.SetWithExpiration(ctx, "session:abc", "xyz", time.Hour)
func (r *Redis) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// SetList pushes a list of string values into a Redis list (LPUSH)
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.SetList(ctx, "mylist", []string{"apple", "banana", "cherry"})
func (r *Redis) SetList(ctx context.Context, key string, values []string) error {
	if len(values) == 0 {
		return nil
	}
	interfaceSlice := make([]interface{}, len(values))
	for i, v := range values {
		interfaceSlice[i] = v
	}
	return r.Client.LPush(ctx, key, interfaceSlice...).Err()
}

// Expire sets an expiration time on an existing key
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.Expire(ctx, "user:123", 30*time.Minute)
func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

// Get retrieves the value of a key
//
// Example:
//
//	ctx := context.Background()
//	value, err := redisClient.Get(ctx, "user:123")
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete removes a key
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.Delete(ctx, "user:123")
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// Exists checks if a key exists
//
// Example:
//
//	ctx := context.Background()
//	exists, err := redisClient.Exists(ctx, "user:123")
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Close shuts down the Redis connection
//
// Example:
//
//	err := redisClient.Close()
func (r *Redis) Close() error {
	return r.Client.Close()
}
