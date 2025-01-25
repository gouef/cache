package cache

import (
	"context"
	"errors"
	"github.com/gouef/standards"
	redisLib "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redisLib.Client
	ctx    context.Context
}

func NewRedisCache(addr, password string, db int) standards.Cache {
	client := redisLib.NewClient(&redisLib.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *Redis) GetItem(key string) standards.CacheItem {
	value, err := c.client.Get(c.ctx, key).Result()
	if err == redisLib.Nil {
		return &RedisItem{key: key, hit: false}
	}
	return &RedisItem{key: key, value: value, hit: true}
}

func (c *Redis) GetItems(keys ...string) []standards.CacheItem {
	var items []standards.CacheItem
	for _, key := range keys {
		items = append(items, c.GetItem(key))
	}
	return items
}

func (c *Redis) HasItem(key string) bool {
	_, err := c.client.Get(c.ctx, key).Result()
	return err != redisLib.Nil
}

func (c *Redis) Clear() error {
	return c.client.FlushAll(c.ctx).Err()
}

func (c *Redis) DeleteItem(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *Redis) DeleteItems(keys ...string) error {
	return c.client.Del(c.ctx, keys...).Err()
}

func (c *Redis) Save(item standards.CacheItem) error {
	rItem, ok := item.(*RedisItem)
	if !ok {
		return errors.New("invalid cache item type")
	}
	return c.client.Set(c.ctx, rItem.GetKey(), rItem.Get(), 0).Err()
}

func (c *Redis) SaveDeferred(item standards.CacheItem) error {
	return c.Save(item)
}

func (c *Redis) Commit() error {
	return nil
}
