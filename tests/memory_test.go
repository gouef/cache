package tests

import (
	"github.com/gouef/cache"
	"github.com/gouef/standards"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemory(t *testing.T) {
	t.Run("Create Memory and basic functions", func(t *testing.T) {
		memory := cache.NewMemory()
		item, err := cache.NewMemoryItem("test").Set("data", standards.KeepTTL)

		assert.Empty(t, memory.GetItems("test"))
		assert.Nil(t, err)
		assert.Nil(t, memory.Save(item))

		item.ExpiresAfter(3 * time.Minute)

		item2, err := cache.NewMemoryItem("test 2").Set("test data", standards.KeepTTL)
		assert.NoError(t, err)
		item3, err := cache.NewMemoryItem("test 3").Set("test data 3", standards.KeepTTL)
		assert.NoError(t, err)

		assert.NoError(t, memory.Save(item2))
		assert.NoError(t, memory.Save(item3))

		assert.True(t, memory.HasItem("test"))
		assert.NotNil(t, memory.GetItem("test"))

		assert.NoError(t, memory.DeleteItems("test", "test 3"))

		assert.Equal(t, 1, len(memory.GetItems("test", "test 2", "test 3", "non-exists")))

		assert.NoError(t, memory.DeleteItem("test 2"))

		assert.False(t, memory.HasItem("test 2"))
		assert.Equal(t, 0, len(memory.GetItems("test", "test 2", "test 3", "non-exists")))

		assert.NoError(t, memory.Save(item3))
		assert.NoError(t, memory.Clear())

		assert.Equal(t, 0, len(memory.GetItems("test", "test 2", "test 3", "non-exists")))

	})
	t.Run("Try save FileItem ", func(t *testing.T) {
		memory := cache.NewMemory()
		item, err := cache.NewFileItem("test").Set("data", standards.KeepTTL)

		assert.Empty(t, memory.GetItems("test"))
		assert.Nil(t, err)
		assert.Error(t, memory.SaveDeferred(item))

		assert.Nil(t, memory.Commit())
	})
}

func TestMemoryItem(t *testing.T) {
	t.Run("Create MemoryItem and basic functions", func(t *testing.T) {
		item := cache.NewMemoryItem("test")

		assert.NotNil(t, item)
		assert.Equal(t, "test", item.GetKey())
		assert.Equal(t, false, item.IsHit())

		c, err := item.Set("data", standards.KeepTTL)

		assert.NoError(t, err)
		assert.Equal(t, "data", c.Get())

		c.ExpiresAfter(3 * time.Minute)

		assert.True(t, c.IsHit())

		oneSecond := time.Now().Add(1 * time.Second)
		c.ExpiresAt(oneSecond)
		time.Sleep(2 * time.Second)

		assert.False(t, c.IsHit())
		assert.Nil(t, c.Get())
	})

	t.Run("Create MemoryItem function and basic functions", func(t *testing.T) {
		item, err := cache.NewMemoryItem("test").Set(func() int {
			return 7 + 4
		}, 0)

		assert.NoError(t, err)

		assert.NotNil(t, item)
		assert.Equal(t, "test", item.GetKey())
		assert.Equal(t, false, item.IsHit())

		c, err := item.Set("data", standards.KeepTTL)

		assert.NoError(t, err)
		assert.Equal(t, "data", c.Get())

		c.ExpiresAfter(3 * time.Minute)

		assert.True(t, c.IsHit())

		oneSecond := time.Now().Add(1 * time.Second)
		c.ExpiresAt(oneSecond)
		time.Sleep(2 * time.Second)

		assert.False(t, c.IsHit())
		assert.Nil(t, c.Get())
	})
}
