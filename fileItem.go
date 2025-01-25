package cache

import (
	"github.com/gouef/standards"
	"time"
)

type FileItem struct {
	Key        string    `json:"key"`
	Value      any       `json:"value"`
	Expiration time.Time `json:"expiration"`
}

func NewFileItem(key string) *FileItem {
	return &FileItem{Key: key}
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
	return i.Expiration.IsZero() || i.Expiration.After(time.Now())
}

func (i *FileItem) Set(value any) (standards.CacheItem, error) {
	i.Value = value
	return i, nil
}

func (i *FileItem) ExpiresAt(expiration time.Time) (standards.CacheItem, error) {
	i.Expiration = expiration
	return i, nil
}

func (i *FileItem) ExpiresAfter(t time.Duration) {
	i.ExpiresAt(time.Now().Add(t))
}
