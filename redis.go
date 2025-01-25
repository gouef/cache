package cache

import (
	"context"
	"errors"
	"github.com/gouef/standards"
	redisLib "github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redisLib.Client
	ctx    context.Context
}

func NewRedis(client *redisLib.Client) standards.Cache {
	return &Redis{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *Redis) GetItem(key string) standards.CacheItem {
	value, err := c.client.Get(c.ctx, key).Result()
	if err == redisLib.Nil {
		return nil
	}
	return &RedisItem{key: key, value: value, hit: true}
}

func (c *Redis) GetItems(keys ...string) []standards.CacheItem {
	var items []standards.CacheItem
	for _, key := range keys {
		item := c.GetItem(key)
		if item != nil && item.Get() != "" {
			items = append(items, c.GetItem(key))
		}
	}
	return items
}

func (c *Redis) HasItem(key string) bool {
	item, err := c.client.Get(c.ctx, key).Result()
	return err != redisLib.Nil && item != ""
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
	return c.client.Set(c.ctx, rItem.GetKey(), rItem.Get(), rItem.expiration.Sub(time.Now())).Err()
}

func (c *Redis) SaveDeferred(item standards.CacheItem) error {
	return c.Save(item)
}

func (c *Redis) Commit() error {
	return nil
}
