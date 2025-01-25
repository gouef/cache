package cache

import (
	"github.com/gouef/standards"
	"time"
)

type RedisItem struct {
	key        string
	value      any
	hit        bool
	expiration time.Time
}

func NewRedisItem(key string) *RedisItem {
	return &RedisItem{
		key: key,
		hit: false,
	}
}

func (r *RedisItem) GetKey() string {
	return r.key
}

func (r *RedisItem) Get() any {
	if r.IsHit() {
		return r.value
	}

	return nil
}

func (r *RedisItem) IsHit() bool {
	return r.hit && (r.expiration.IsZero() || r.expiration.After(time.Now()))
}

func (r *RedisItem) Set(value any) (standards.CacheItem, error) {
	r.value = value
	r.hit = true
	return r, nil
}

func (r *RedisItem) ExpiresAt(expiration time.Time) (standards.CacheItem, error) {
	r.expiration = expiration
	return r, nil
}

func (r *RedisItem) ExpiresAfter(t time.Duration) {
	r.ExpiresAt(time.Now().Add(t))
}
