package cache

import (
	"github.com/gouef/standards"
	"time"
)

type FileItem struct {
	Key        string    `json:"key"`
	Value      any       `json:"value"`
	Expiration time.Time `json:"expiration"`
	KeepTTL    bool
}

func NewFileItem(key string) *FileItem {
	return &FileItem{Key: key, KeepTTL: false}
}

func (i *FileItem) GetKey() string {
	return i.Key
}

func (i *FileItem) Get() any {
	if i.IsHit() {
		return i.Value
	}
	return nil
}

func (i *FileItem) IsHit() bool {

	return i.KeepTTL || i.Expiration.IsZero() || i.Expiration.After(time.Now())
}

func (i *FileItem) Set(value any, ttl time.Duration) (standards.CacheItem, error) {
	i.Value = value
	if ttl == KeepTTL {
		i.KeepTTL = true
	}
	return i, nil
}

func (i *FileItem) ExpiresAt(expiration time.Time) (standards.CacheItem, error) {
	i.Expiration = expiration
	i.KeepTTL = false
	return i, nil
}

func (i *FileItem) ExpiresAfter(t time.Duration) {
	if t == KeepTTL {
		i.KeepTTL = true
	} else {
		i.ExpiresAt(time.Now().Add(t))
	}
}
