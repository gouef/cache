package cache

import (
	"github.com/gouef/standards"
	"sync"
	"time"
)

type MemoryItem struct {
	key        string
	value      any
	expiration time.Time
	KeepTTL    bool
	hit        bool
	mu         sync.RWMutex
}

func NewMemoryItem(key string) *MemoryItem {
	return &MemoryItem{
		key: key,
		hit: false,
	}
}

func (m *MemoryItem) GetKey() string {
	return m.key
}

func (m *MemoryItem) Get() any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.IsHit() {
		return m.value
	}
	return nil
}

func (m *MemoryItem) IsHit() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.KeepTTL || (m.hit && (m.expiration.IsZero() || m.expiration.After(time.Now())))
}

func (m *MemoryItem) Set(value any, ttl time.Duration) (standards.CacheItem, error) {

	m.mu.Lock()
	defer m.mu.Unlock()
	m.value = value
	m.ExpiresAfter(ttl)
	return m, nil
}

func (m *MemoryItem) ExpiresAt(expiration time.Time) (standards.CacheItem, error) {
	m.expiration = expiration
	m.KeepTTL = false
	return m, nil
}

func (m *MemoryItem) ExpiresAfter(t time.Duration) {
	if t == KeepTTL {
		m.KeepTTL = true
		m.hit = true
	} else {
		m.ExpiresAt(time.Now().Add(t))
	}
}
