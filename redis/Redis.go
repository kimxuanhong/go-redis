package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis defines the interface for interacting with Client
type Redis interface {
	Set(ctx context.Context, key string, value interface{}) error
	SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetList(ctx context.Context, key string, values []string) error
	GetList(ctx context.Context, key string, start, stop int64) ([]string, error)
	LPop(ctx context.Context, key string) (string, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Get(ctx context.Context, key string) (string, error)
	Increment(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	Close() error
}

// Client wraps a Client client
type Client struct {
	*redis.Client
}

// NewRedis initializes a new Client client based on the given configuration
//
// Example:
//
//	cfg := &Config{}
//	redisClient, err := NewRedis(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewRedis(cfg *Config) (Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Client")
	return &Client{client}, nil
}

// Set stores a key-value pair without expiration
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.Set(ctx, "user:123", "John Doe")
func (r *Client) Set(ctx context.Context, key string, value interface{}) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}

// SetWithExpiration stores a key-value pair with expiration
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.SetWithExpiration(ctx, "session:abc", "xyz", time.Hour)
func (r *Client) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// SetList pushes a list of string values into a Client list (LPUSH)
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.SetList(ctx, "mylist", []string{"apple", "banana", "cherry"})
func (r *Client) SetList(ctx context.Context, key string, values []string) error {
	if len(values) == 0 {
		return nil
	}
	interfaceSlice := make([]interface{}, len(values))
	for i, v := range values {
		interfaceSlice[i] = v
	}
	return r.Client.LPush(ctx, key, interfaceSlice...).Err()
}

// GetList retrieves a range of elements from a Redis list using LRANGE.
//
// Example:
//
//	values, err := redisClient.GetList(ctx, "mylist", 0, -1)
func (r *Client) GetList(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.LRange(ctx, key, start, stop).Result()
}

// LPop removes and returns the first element from a Redis list.
//
// Example:
//
//	val, err := redisClient.LPop(ctx, "mylist")
func (r *Client) LPop(ctx context.Context, key string) (string, error) {
	return r.Client.LPop(ctx, key).Result()
}

// Expire sets an expiration time on an existing key
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.Expire(ctx, "user:123", 30*time.Minute)
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

// SetNX sets a key with a value only if the key does not already exist.
// It's commonly used for implementing distributed locks.
//
// Example:
//
//	ok, err := redisClient.SetNX(ctx, "lock:job1", "locked", 10*time.Second)
func (r *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}

// TTL returns the time-to-live for a key. If the key has no expiration, it returns -1.
//
// Example:
//
//	ttl, err := redisClient.TTL(ctx, "session:abc")
func (r *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

// Get retrieves the value of a key
//
// Example:
//
//	ctx := context.Background()
//	value, err := redisClient.Get(ctx, "user:123")
func (r *Client) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Increment increments the integer value of a key by 1.
//
// Example:
//
//	ctx := context.Background()
//	newVal, err := redisClient.Increment(ctx, "counter:pageviews")
func (r *Client) Increment(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// Delete removes a key
//
// Example:
//
//	ctx := context.Background()
//	err := redisClient.Delete(ctx, "user:123")
func (r *Client) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// Exists checks if a key exists
//
// Example:
//
//	ctx := context.Background()
//	exists, err := redisClient.Exists(ctx, "user:123")
func (r *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetJSON serializes a value to JSON and stores it in Redis with expiration.
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	err := redisClient.SetJSON(ctx, "user:123", User{Name: "Alice", Age: 30}, time.Hour)
func (r *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, data, expiration).Err()
}

// GetJSON retrieves a value from Redis and unmarshals it into the provided destination.
//
// Example:
//
//	var user User
//	err := redisClient.GetJSON(ctx, "user:123", &user)
func (r *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Close shuts down the Client connection
//
// Example:
//
//	err := redisClient.Close()
func (r *Client) Close() error {
	return r.Client.Close()
}
