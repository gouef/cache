package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/gouef/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	t.Run("Redis with basic functions", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		mock.ExpectGet("test").SetErr(errors.New("not found"))

		r := cache.NewRedis(db)
		item, err := cache.NewRedisItem("test").Set("data")

		assert.Empty(t, r.GetItems("test"))
		assert.Nil(t, err)
		assert.False(t, r.HasItem("test"))

		mock.ExpectSet("test", "data", 0).SetVal("data")
		assert.Nil(t, r.Save(item))

		item.ExpiresAfter(3 * time.Minute)

		item2, err := cache.NewRedisItem("test 2").Set("test data")
		assert.NoError(t, err)
		item3, err := cache.NewRedisItem("test 3").Set("test data 3")
		assert.NoError(t, err)

		mock.ExpectSet("test 2", "test data", 0).SetVal("test data")
		assert.NoError(t, r.Save(item2))
		mock.ExpectSet("test 3", "test data 3", 0).SetVal("test data 3")
		assert.NoError(t, r.Save(item3))

		mock.ExpectGet("test").SetVal("data")
		assert.True(t, r.HasItem("test"))
		assert.NotNil(t, r.GetItem("test"))

		mock.ExpectDel("test", "test 3").SetVal(0)
		assert.NoError(t, r.DeleteItems("test", "test 3"))

		mock.ExpectGet("test").SetVal("")
		mock.ExpectGet("test 2").SetVal("test data")
		mock.ExpectGet("test 3").SetVal("")
		mock.ExpectGet("non-exists").SetVal("")
		assert.Equal(t, 1, len(r.GetItems("test", "test 2", "test 3", "non-exists")))

		mock.ExpectDel("test 2").SetVal(0)
		assert.NoError(t, r.DeleteItem("test 2"))

		mock.ExpectGet("test 2").SetVal("")
		assert.False(t, r.HasItem("test 2"))

		mock.ExpectGet("test").SetVal("")
		mock.ExpectGet("test 2").SetVal("")
		mock.ExpectGet("test 3").SetVal("")
		mock.ExpectGet("non-exists").SetVal("")
		assert.Equal(t, 0, len(r.GetItems("test", "test 2", "test 3", "non-exists")))

		mock.ExpectSet("test 3", "test data 3", 0).SetVal("test data 3")
		assert.NoError(t, r.SaveDeferred(item3))

		mock.ExpectFlushAll().SetVal("")
		assert.NoError(t, r.Clear())

		mock.ExpectGet("test").SetVal("")
		mock.ExpectGet("test 2").SetVal("")
		mock.ExpectGet("test 3").SetVal("")
		mock.ExpectGet("non-exists").SetVal("")
		assert.Equal(t, 0, len(r.GetItems("test", "test 2", "test 3", "non-exists")))

		assert.Nil(t, r.Commit())

		mock.ExpectGet("rediLibNil").RedisNil()
		assert.Nil(t, r.GetItem("rediLibNil"))

		mItem, err := cache.NewMemoryItem("mTest").Set("data memory")
		assert.Nil(t, err)

		assert.Error(t, r.Save(mItem))
	})
}

func TestRedisItem(t *testing.T) {
	t.Run("RedisItem with basic functions", func(t *testing.T) {
		item := cache.NewRedisItem("test")

		assert.NotNil(t, item)
		assert.Equal(t, "test", item.GetKey())
		assert.Equal(t, false, item.IsHit())

		c, err := item.Set("data")

		assert.NoError(t, err)
		assert.Equal(t, "data", c.Get())

		c.ExpiresAfter(3 * time.Minute)

		assert.True(t, c.IsHit())

		c.ExpiresAt(time.Now().Add(1 * time.Second))
		time.Sleep(2 * time.Second)

		assert.False(t, c.IsHit())
		assert.Nil(t, c.Get())
	})
}

var ctx = context.TODO()

func NewsInfoForCache(redisDB *redis.Client, newsID int) (info string, err error) {
	cacheKey := fmt.Sprintf("news_redis_cache_%d", newsID)
	info, err = redisDB.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// info, err = call api()
		info = "test"
		err = redisDB.Set(ctx, cacheKey, info, 30*time.Minute).Err()
	}
	return
}

func TestNewsInfoForCache(t *testing.T) {
	db, mock := redismock.NewClientMock()

	newsID := 123456789
	key := fmt.Sprintf("news_redis_cache_%d", newsID)

	// mock ignoring `call api()`

	mock.ExpectGet(key).RedisNil()
	mock.Regexp().ExpectSet(key, `[a-z]+`, 30*time.Minute).SetErr(errors.New("FAIL"))

	_, err := NewsInfoForCache(db, newsID)
	if err == nil || err.Error() != "FAIL" {
		t.Error("wrong error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
