# Memory cache
This package provides an in-memory cache system for storing and retrieving data. It supports basic operations such as saving items, retrieving items, deleting items, and handling expiration.

## Types

### Memory
The main structure for working with an in-memory cache.

#### Properties:
- `items`: A map that holds the cache items.
- `mu`: A mutex to ensure thread-safety when accessing and modifying the cache items.

#### Functions:
- `NewMemory()`: Creates a new instance of Memory.

```go
cache := NewMemory()
```

- `GetItem(key string) standards.CacheItem`: Retrieves a cache item by its key. If the item doesn't exist, it returns nil.

```go
item := cache.GetItem("some-key")
```

- `GetItems(keys ...string) []standards.CacheItem`: Retrieves multiple cache items by a list of keys.

```go
items := cache.GetItems("key1", "key2")
```

- `HasItem(key string) bool`: Checks if a cache item with the given key exists and is still valid.

```go
exists := cache.HasItem("some-key")
```

- `Clear() error`: Clears all items in the cache.

```go
err := cache.Clear()
```

- `DeleteItem(key string) error`: Deletes a cache item by its key.

```go
err := cache.DeleteItem("some-key")
```

- `DeleteItems(keys ...string) error`: Deletes multiple cache items by a list of keys.

```go
err := cache.DeleteItems("key1", "key2")
```

- `Save(item standards.CacheItem) error`: Saves a cache item.

```go
err := cache.Save(item)
```

- `SaveDeferred(item standards.CacheItem) error`: Saves a cache item (same as Save).

```go
err := cache.SaveDeferred(item)
```

- `Commit() error`: Not currently implemented, but it's available for future operations.

### MemoryItem
A structure representing an individual cache item in the memory-based cache.

#### Properties:
- `key`: The cache item's key.
- `value`: The cache item's value, which can be of any type (any).
- `expiration`: The expiration time of the cache item.
- `hit`: A flag indicating whether the item is valid and has been accessed.

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
	"github.com/gouef/cache"
)

func main() {
	// Create a new in-memory cache instance
	c := cache.NewMemory()

	// Create a new cache item
	item := cache.NewMemoryItem("user123")
	item.Set("some data")
	item.ExpiresAfter(5 * time.Minute)

	// Save the item to cache
	err := c.Save(item)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the item from cache
	loadedItem := c.GetItem("user123")
	if loadedItem != nil {
		fmt.Println("Loaded item:", loadedItem.Get())
	}

	// Delete the item
	err = c.DeleteItem("user123")
	if err != nil {
		log.Fatal(err)
	}
}
```