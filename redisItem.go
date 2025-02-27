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
	KeepTTL    bool
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
	return r.KeepTTL || (r.hit && (r.expiration.IsZero() || r.expiration.After(time.Now())))
}

func (r *RedisItem) Set(value any, ttl time.Duration) (standards.CacheItem, error) {
	r.value = value
	r.hit = true
	if ttl == KeepTTL {
		r.KeepTTL = true
	}
	return r, nil
}

func (r *RedisItem) ExpiresAt(expiration time.Time) (standards.CacheItem, error) {
	r.expiration = expiration
	r.KeepTTL = false
	return r, nil
}

func (r *RedisItem) ExpiresAfter(t time.Duration) {
	if t == KeepTTL {
		r.KeepTTL = true
	} else {
		r.ExpiresAt(time.Now().Add(t))
	}
}
