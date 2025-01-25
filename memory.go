package cache

import (
	"errors"
	"github.com/gouef/standards"
	"sync"
)

type Memory struct {
	items map[string]*MemoryItem
	mu    sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		items: make(map[string]*MemoryItem),
	}
}

func (c *Memory) GetItem(key string) standards.CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.items[key]
	if !exists {
		return nil
	}
	return item
}

func (c *Memory) GetItems(keys ...string) []standards.CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []standards.CacheItem
	for _, key := range keys {
		item := c.GetItem(key)
		if item != nil {
			result = append(result, c.GetItem(key))
		}
	}
	return result
}

func (c *Memory) HasItem(key string) bool {
	item := c.GetItem(key)
	if item == nil {
		return false
	}
	return item.IsHit()
}

func (c *Memory) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*MemoryItem)
	return nil
}

func (c *Memory) DeleteItem(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *Memory) DeleteItems(keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.items, key)
	}
	return nil
}

func (c *Memory) Save(item standards.CacheItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	mItem, ok := item.(*MemoryItem)
	if !ok {
		return errors.New("invalid cache item type")
	}
	c.items[mItem.GetKey()] = mItem
	return nil
}

func (c *Memory) SaveDeferred(item standards.CacheItem) error {
	return c.Save(item)
}

func (c *Memory) Commit() error {
	return nil
}
