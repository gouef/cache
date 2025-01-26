package tests

import (
	"github.com/go-redis/redismock/v9"
	"github.com/gouef/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("Storage", func(t *testing.T) {
		storage := cache.NewStorage()
		assert.NotNil(t, storage)

		memory, err := storage.AddMemory("memory")
		assert.NotNil(t, memory)
		assert.Nil(t, err)

		getMemory, exists := storage.Get("memory")

		assert.NotNil(t, getMemory)
		assert.True(t, exists)

		memoryAlready, err := storage.AddMemory("memory")
		assert.NotNil(t, memoryAlready)
		assert.Error(t, err)

		file, err := storage.AddFile("file", "cache")
		assert.NotNil(t, file)
		assert.Nil(t, err)

		fileErr, err := storage.AddFile("file", "/non-exists")
		assert.Nil(t, fileErr)
		assert.Error(t, err)

		db, _ := redismock.NewClientMock()
		redis, err := storage.AddRedis("redis", db)
		assert.NotNil(t, redis)
		assert.Nil(t, err)
	})
}
