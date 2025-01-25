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

type File struct {
	Dir string
	Mu  sync.RWMutex
}

const FILE_EXTENSION = ".cache"

// NewFile create new instance of File and check if directory exists.
func NewFile(dir string) (standards.Cache, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return &File{Dir: dir}, nil
}

func (c *File) GetItem(key string) standards.CacheItem {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	filePath := c.getFilePath(key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	var item *FileItem
	if err := json.Unmarshal(data, &item); err != nil {
		_ = os.Remove(filePath)
		return nil
	}

	if item.Expiration.Before(time.Now()) && !item.Expiration.IsZero() {
		_ = os.Remove(filePath)
		return nil
	}

	return item
}

func (c *File) GetItems(keys ...string) []standards.CacheItem {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	var items []standards.CacheItem
	for _, key := range keys {
		item := c.GetItem(key)
		if item != nil {
			items = append(items, c.GetItem(key))
		}
	}
	return items
}

func (c *File) HasItem(key string) bool {
	item := c.GetItem(key)

	if item == nil {
		return false
	}

	return item.IsHit()
}

func (c *File) Clear() error {
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

func (c *File) DeleteItem(key string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	return os.Remove(c.getFilePath(key))
}

func (c *File) DeleteItems(keys ...string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	for _, key := range keys {
		_ = os.Remove(c.getFilePath(key))
	}
	return nil
}

func (c *File) Save(item standards.CacheItem) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	fItem, ok := item.(*FileItem)
	if !ok {
		return errors.New("invalid cache item type")
	}

	data, err := json.Marshal(fItem)
	if err != nil {
		return err
	}

	return os.WriteFile(c.getFilePath(fItem.Key), data, 0644)
}

func (c *File) SaveDeferred(item standards.CacheItem) error {
	return c.Save(item)
}

func (c *File) Commit() error {
	return nil
}

func (c *File) getFilePath(key string) string {
	return filepath.Join(c.Dir, key+FILE_EXTENSION)
}
