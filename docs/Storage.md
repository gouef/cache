# Storage
It's wrapper of `Cache` instances.

## Functions:
- `NewStorage() *Storage`: create new instance of `Storage`
- `GetStorage() *Storage`: get created instance of `Storage` (global usages).
- `Add(name string, cache standards.Cache) (standards.Cache, error)`: add cache instance to list
- `Get(name string) (cache standards.Cache, exists bool)`: return cache instance
- `AddFile(name, dir string) (standards.Cache, error)`: create File cache instance and add it to list
- `AddMemory(name string) (standards.Cache, error)`: create Memory cache instance and add it to list
- `AddRedis(name string, client *redisLib.Client) (standards.Cache, error)`: create Redis cache instance and add it to list

## Example usage


```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/gouef/cache"
)

func main() {
	storage := cache.NewStorage()
	
	createdUsersCache, _ := storage.AddMemory("users")

	// Create a new cache item
	item := cache.NewMemoryItem("user123")
	item.Set("some data")
	item.ExpiresAfter(5 * time.Minute)

	// Save the item to cache
	err := createdUsersCache.Save(item)
	if err != nil {
		log.Fatal(err)
	}


	// Returns Memory cache instance
	usersCache, _ := storage.Get("users")

	// Retrieve the item from cache
	loadedItem := usersCache.GetItem("user123")
	if loadedItem != nil {
		fmt.Println("Loaded item:", loadedItem.Get())
	}

	// Delete the item
	err = usersCache.DeleteItem("user123")
	if err != nil {
		log.Fatal(err)
	}
}

func write() {
    storage := cache.GetStorage()
	// ...
}
```