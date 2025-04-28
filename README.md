# go-redis

A simple and lightweight Redis client library for Go applications.

## Features

- 🚀 Simple and clean interface
- 🔒 Thread-safe Redis client
- ⏱️ Context support for all operations
- 🔄 Connection testing on initialization
- 📦 Minimal dependencies
- 🛠️ Basic Redis operations:
  - Set/Get key-value pairs with expiration
  - Delete keys
  - Check key existence

## Installation

```bash
go get github.com/kimxuanhong/go-redis
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kimxuanhong/go-redis/pkg"
)

func main() {
	// Initialize Redis client with default configuration
	redisClient, err := pkg.NewRedis(pkg.NewRedisConfig())
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	ctx := context.Background()

	// Set a key with expiration
	err = redisClient.Set(ctx, "user:1", "John Doe", time.Hour)
	if err != nil {
		panic(err)
	}

	// Get a value
	value, err := redisClient.Get(ctx, "user:1")
	if err != nil {
		panic(err)
	}
	fmt.Println("User:", value)

	// Check if key exists
	exists, err := redisClient.Exists(ctx, "user:1")
	if err != nil {
		panic(err)
	}
	fmt.Println("Key exists:", exists)

	// Delete a key
	err = redisClient.Delete(ctx, "user:1")
	if err != nil {
		panic(err)
	}
}
```

## Configuration

The Redis client can be configured using environment variables or by creating a custom `RedisConfig`:

### Environment Variables

```bash
REDIS_HOST=localhost     # Default: localhost
REDIS_PORT=6379         # Default: 6379
REDIS_PASSWORD=         # Default: empty
REDIS_DB=0              # Default: 0
```

### Custom Configuration

```go
config := &pkg.RedisConfig{
    Host:     "localhost",
    Port:     "6379",
    Password: "your-password",
    DB:       0,
}

redisClient, err := pkg.NewRedis(config)
```

## Error Handling

All operations return an error that should be checked:

```go
value, err := redisClient.Get(ctx, "key")
if err != nil {
    if err == redis.Nil {
        // Key does not exist
    } else {
        // Other error occurred
    }
}
```

## Best Practices

1. Always close the client when done:
```go
defer redisClient.Close()
```

2. Use context for timeouts and cancellation:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

3. Handle connection errors gracefully:
```go
redisClient, err := pkg.NewRedis(config)
if err != nil {
    // Handle connection error
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 