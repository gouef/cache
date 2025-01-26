# File cache
This package provides a simple file-based cache system for storing and retrieving data from files in a specified directory. It supports basic operations such as saving items, retrieving items, deleting items, and handling expiration.

## Types

### File
The main structure for working with the file-based cache.

#### Properties:
- `Dir`: The directory path where cache files will be stored.
- `Mu`: A lock to ensure thread-safety when accessing the files.

#### Functions:
- `NewFile(dir string)`: Creates a new instance of File and checks if the directory exists. If it doesn't, it will create it.

```go
cache, err := NewFile("/path/to/cache")
```

- `GetItem(key string) standards.CacheItem`: Retrieves a cache item by its key. If the item doesn't exist or has expired, it returns nil.

```go
item := cache.GetItem("some-key")
```

- `GetItems(keys ...string) []standards.CacheItem`: Retrieves multiple items by a list of keys.

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

### FileItem
A structure representing an individual cache item.

#### Properties:
- `Key`: The cache item's key.
- `Value`: The cache item's value, which can be of any type (any).
- `Expiration`: The expiration time of the cache item.

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

- `Set(value any)`: Sets the value of the item.

```go
item.Set(newValue)
```

- `ExpiresAt(expiration time.Time): Sets` the explicit expiration time for the item.

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
	// Create a new cache instance
	c, err := cache.NewFile("/path/to/cache")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new cache item
	item := cache.NewFileItem("user123")
	item.Set("some data")
	item.ExpiresAfter(5 * time.Minute)

	// Save the item to cache
	err = c.Save(item)
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

### FileSimple
The main simplified structure for working with the file-based cache.

#### Properties:
- `Dir`: The directory path where cache files will be stored.
- `Mu`: A lock to ensure thread-safety when accessing the files.
- `AllowDefaultNil`: if `true` it will return `nil` for case when cache item not exists.

#### Functions:

- `NewFileSimple(dir string) (*FileSimple, error)`: create FileSimple instance (allowDefaultNil `false`)
- `NewFileSimpleWithDefaultNil(dir string, allowDefaultNil bool) (*FileSimple, error)`: create FileSimple instance
- `Get(key string, defaultValue any) any`: Returns a value from the cache.
- `GetMultiply(keys []string, defaultValue any) []any`: Returns a list of cache items.
- `Has(key string) bool`: Determines whether an item is present in the cache.
- `Clear() error`: Deletes all cache's keys.
- `Delete(key string) error`: Remove an item from the cache.
- `DeleteMultiply(keys ...string) error`: Removes multiple items in a single operation.
- `Set(key string, item any) error`: Persists a cache item.
- `SetMultiply(values map[string]any, ttl time.Duration) error`: Persists a cache items.

