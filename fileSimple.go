package cache

import (
	"encoding/json"
	"errors"
	"github.com/gouef/standards"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileSimple struct {
	Dir             string
	Mu              sync.RWMutex
	AllowDefaultNil bool
	error           error
}

// NewFileSimple create FileSimple instance with not allowed default value nil
func NewFileSimple(dir string) (*FileSimple, error) {
	return NewFileSimpleWithDefaultNil(dir, false)
}

// NewFileSimpleWithDefaultNil create FileSimple instance
func NewFileSimpleWithDefaultNil(dir string, allowDefaultNil bool) (*FileSimple, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return &FileSimple{
		Dir:             dir,
		AllowDefaultNil: allowDefaultNil,
	}, nil
}

// Get Returns a value from the cache.
func (c *FileSimple) Get(key string, defaultValue any) any {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	filePath := c.getFilePath(key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return defaultValue
	}

	var item *FileItem
	if err := json.Unmarshal(data, &item); err != nil {
		_ = os.Remove(filePath)
		return defaultValue
	}

	if time.Now().After(item.Expiration) && !item.Expiration.IsZero() {
		_ = os.Remove(filePath)
		return defaultValue
	}

	return item.Value
}

// GetMultiply Returns a list of cache items.
func (c *FileSimple) GetMultiply(keys []string, defaultValue any) []any {
	result := []any{}

	for _, key := range keys {
		item := c.Get(key, defaultValue)

		if (c.AllowDefaultNil && item == nil) || item != nil {
			result = append(result, item)
		}
	}

	return result
}

// Has Determines whether an item is present in the cache.
func (c *FileSimple) Has(key string) bool {
	item := c.Get(key, nil)

	if (c.AllowDefaultNil && item == nil) || item != nil {
		return true
	}

	return false
}

// Clear Deletes all cache's keys.
func (c *FileSimple) Clear() error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	files, err := os.ReadDir(c.Dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := strings.Replace(file.Name(), FILE_EXTENSION, "", 1)
		err := os.Remove(c.getFilePath(name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Remove an item from the cache.
func (c *FileSimple) Delete(key string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	return os.Remove(c.getFilePath(key))
}

// DeleteMultiply Removes multiple items in a single operation.
func (c *FileSimple) DeleteMultiply(keys ...string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	for _, key := range keys {
		_ = os.Remove(c.getFilePath(key))
	}
	return nil
}

// Set Persists a cache item.
func (c *FileSimple) Set(key string, item any) error {
	fItem, err := c.getFileItem(key, item)

	data, err := json.Marshal(fItem)
	if err != nil {
		return err
	}

	return os.WriteFile(c.getFilePath(fItem.GetKey()), data, 0644)
}

// SetMultiply Persists a cache items.
func (c *FileSimple) SetMultiply(values map[string]any, ttl time.Duration) error {
	for key, value := range values {
		item, err := c.getFileItem(key, value)
		item.ExpiresAfter(ttl)

		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		err = os.WriteFile(c.getFilePath(item.GetKey()), data, 0644)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *FileSimple) getFileItem(key string, value any) (standards.CacheItem, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	item, err := NewFileItem(key).Set(value)
	if err != nil {
		return nil, errors.New("invalid cache item type")
	}

	return item, nil
}

func (c *FileSimple) getFilePath(key string) string {
	return filepath.Join(c.Dir, key+FILE_EXTENSION)
}
