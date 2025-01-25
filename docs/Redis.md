# Redis cache
This package provides a Redis-based cache system for storing and retrieving data. It supports basic cache operations such as saving, retrieving, deleting items, and handling expiration.

## Types

### Redis

#### Properties:
- `client`: A Redis client used to communicate with the Redis server.
- `ctx`: The context used for Redis operations.

#### Functions:

- `NewRedis(client *redisLib.Client) standards.Cache`: Creates a new Redis cache instance using an existing Redis client.

```go
redisCache := NewRedis(client)
```

- `GetItem(key string) standards.CacheItem`: Retrieves a cache item by its key from Redis. If the item doesn't exist, it returns nil.

```go
item := redisCache.GetItem("some-key")
```

- `GetItems(keys ...string) []standards.CacheItem`: Retrieves multiple cache items by a list of keys from Redis.

```go
items := redisCache.GetItems("key1", "key2")
```

- `HasItem(key string) bool`: Checks if a cache item with the given key exists in Redis.

```go
exists := redisCache.HasItem("some-key")
```

- `Clear() error`: Clears all items in the Redis cache.

```go
err := redisCache.Clear()
```

- `DeleteItem(key string) error`: Deletes a cache item by its key from Redis.

```go
err := redisCache.DeleteItem("some-key")
```

- `DeleteItems(keys ...string) error`: Deletes multiple cache items by a list of keys from Redis.

```go
err := redisCache.DeleteItems("key1", "key2")
```

- `Save(item standards.CacheItem) error`: Saves a cache item to Redis. The item is stored with an expiration time.

```go
err := redisCache.Save(item)
```

- `SaveDeferred(item standards.CacheItem) error`: Saves a cache item to Redis (same as Save).

```go
err := redisCache.SaveDeferred(item)
```

- `Commit() error`: Not currently implemented, but it's available for future operations.

### RedisItem
Represents a cache item that is stored in Redis.

#### Properties:
- `key`: The key of the cache item.
- `value`: The value of the cache item (can be of any type).
- `hit`: A flag indicating whether the item is valid and has been accessed.
- `expiration`: The expiration time of the cache item.

#### Functions:
- `GetKey()`: Returns the key of the item.

```go
key := item.GetKey()
```

- `Get()`: Returns the value of the item if it is still valid, otherwise nil.

```go
value := item.Get()
```

- `IsHit()`: Checks if the item is still valid.

```go
isValid := item.IsHit()
```

- `Set(value any)`: Sets the value of the item and marks it as "hit".

```go
item.Set(newValue)
```

- `ExpiresAt(expiration time.Time)`: Sets the explicit expiration time for the item.

```go
item.ExpiresAt(time.Now().Add(time.Hour))
```

- `ExpiresAfter(t time.Duration)`: Sets the expiration time for the item to be after a specified duration.

```go
item.ExpiresAfter(10 * time.Minute)
```

## Example Usage

```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/gouef/cache"
)

func main() {
	// Set up Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Create a new Redis cache instance
	redisCache := cache.NewRedis(client)

	// Create a new cache item
	item := cache.NewRedisItem("user123")
	item.Set("some data")
	item.ExpiresAfter(5 * time.Minute)

	// Save the item to Redis
	err := redisCache.Save(item)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the item from Redis
	loadedItem := redisCache.GetItem("user123")
	if loadedItem != nil {
		fmt.Println("Loaded item:", loadedItem.Get())
	}

	// Delete the item
	err = redisCache.DeleteItem("user123")
	if err != nil {
		log.Fatal(err)
	}
}
```