package cache

import (
	"github.com/gouef/standards"
	"time"
)

type RedisItem struct {
	key   string
	value any
	hit   bool
}

func (r *RedisItem) GetKey() string {
	return r.key
}

func (r *RedisItem) Get() any {
	return r.value
}

func (r *RedisItem) IsHit() bool {
	return r.hit
}

func (r *RedisItem) Set(value any) (standards.CacheItem, error) {
	r.value = value
	r.hit = true
	return r, nil
}

func (r *RedisItem) ExpiresAt(expiration time.Time) (standards.CacheItem, error) {
	return r, nil
}

func (r *RedisItem) ExpiresAfter(t time.Duration) {
}
