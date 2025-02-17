package cache

import (
	"errors"
	"fmt"
	"github.com/gouef/standards"
	redisLib "github.com/redis/go-redis/v9"
	"sync"
)

var (
	mu      sync.RWMutex
	storage *Storage
)

type Storage struct {
	Storages map[string]standards.Cache
}

// NewStorage create new instance of Storage
func NewStorage() *Storage {
	mu.Lock()
	defer mu.Unlock()
	s := &Storage{
		Storages: make(map[string]standards.Cache),
	}
	storage = s
	return s
}

// GetStorage get storage (for global usages)
func GetStorage() *Storage {
	mu.RLock()
	defer mu.RUnlock()
	return storage
}

// Add add cache instance to list
func (s *Storage) Add(name string, cache standards.Cache) (standards.Cache, error) {
	v, exists := s.Get(name)

	if exists {
		return v, errors.New(fmt.Sprintf("Storage with name \"%s\" already exists.", name))
	}

	s.Storages[name] = cache
	return cache, nil
}

// Get return cache instance
func (s *Storage) Get(name string) (cache standards.Cache, exists bool) {
	cache, exists = s.Storages[name]
	return
}

// AddFile create File cache instance and add it to list
func (s *Storage) AddFile(name, dir string) (standards.Cache, error) {
	fileCache, err := NewFile(dir)
	if err != nil {
		return nil, err
	}

	return s.Add(name, fileCache)
}

// AddMemory create Memory cache instance and add it to list
func (s *Storage) AddMemory(name string) (standards.Cache, error) {
	memoryCache := NewMemory()
	return s.Add(name, memoryCache)
}

// AddRedis create Redis cache instance and add it to list
func (s *Storage) AddRedis(name string, client *redisLib.Client) (standards.Cache, error) {
	redisCache := NewRedis(client)
	return s.Add(name, redisCache)
}
