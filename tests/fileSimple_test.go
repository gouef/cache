package tests

import (
	"github.com/gouef/cache"
	"github.com/gouef/standards"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestFileSimple(t *testing.T) {

	t.Run("Clear", func(t *testing.T) {

		dir := setupCacheDir(t)
		c, err := cache.NewFileSimple(dir)
		assert.NoError(t, err, "Failed to create file cache")

		items := []string{"item1", "item2", "item3"}
		for _, key := range items {
			err := c.Set(key, "key"+"_value", standards.KeepTTL)
			assert.NoError(t, err, "Failed to save cache item")
		}

		err = c.Clear()
		assert.NoError(t, err, "Failed to clear cache")

		for _, key := range items {
			cachedItem := c.Get(key, nil)
			assert.Nil(t, cachedItem)
		}
	})

	t.Run("SetMultiply error", func(t *testing.T) {
		dir := setupCacheDir(t)
		c, err := cache.NewFileSimple(dir)
		assert.NoError(t, err, "Failed to create file cache")

		c.Dir = "/non-exists"
		oneMinute := 1 * time.Minute
		assert.Error(t, c.SetMultiply(map[string]any{"item1": "data", "item2": "data2"}, oneMinute))
	})
}

func TestFileSimple_SaveAndGetItem(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFileSimple(dir)
	assert.NoError(t, err, "Failed to create cache directory")

	err = c.Set("example", "Hello, world!", standards.KeepTTL)

	assert.NoError(t, err, "Failed to save cache item")

	cachedItem := c.Get("example", nil)
	assert.NotNil(t, cachedItem)

	if cachedItem != nil {
		assert.True(t, c.Has("example"), "Expected cache item to be a hit")
		assert.Equal(t, "Hello, world!", cachedItem, "Cache item value mismatch")
	}
}

func TestFileSimple_SaveAndGetItems(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFileSimple(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := cache.NewFileItem("example")
	_, err = item.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item.ExpiresAfter(5 * time.Minute)
	err = c.Set(item.GetKey(), item.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item2 := cache.NewFileItem("example2")
	_, err = item2.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item2.ExpiresAfter(5 * time.Minute)

	item3 := cache.NewFileItem("example3")
	_, err = item3.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item3.ExpiresAfter(1 * time.Second)

	err = c.Set(item2.GetKey(), item2.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Set(item3.GetKey(), item3.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")

	time.Sleep(1 * time.Second)

	cachedItems := c.GetMultiply([]string{"example", "example2"}, nil)
	assert.NotEmpty(t, cachedItems)

	assert.Equal(t, 2, len(cachedItems))

	// Invalid json
	mu := sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()

	data := []byte("{{,#)")

	os.WriteFile(filepath.Join(dir, item.GetKey()+".cache"), data, 0644)

	invalidItem := c.Get(item.GetKey(), nil)

	assert.Nil(t, invalidItem)

}

func TestFileSimple_SetMultiply(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFileSimple(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := cache.NewFileItem("example")
	_, err = item.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item.ExpiresAfter(5 * time.Minute)
	err = c.Set(item.GetKey(), item.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item2 := cache.NewFileItem("example2")
	_, err = item2.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item2.ExpiresAfter(5 * time.Minute)

	item3 := cache.NewFileItem("example3")
	_, err = item3.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item3.ExpiresAfter(1 * time.Second)

	err = c.Set(item2.GetKey(), item2.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Set(item3.GetKey(), item3.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")

	data := map[string]any{
		item.GetKey():  item.GetKey(),
		item2.GetKey(): item2.GetKey(),
		item3.GetKey(): item3.GetKey(),
	}
	fiveMinute := 5 * time.Minute
	err = c.SetMultiply(data, fiveMinute)

	time.Sleep(1 * time.Second)

	cachedItems := c.GetMultiply([]string{"example", "example2"}, nil)
	assert.NotEmpty(t, cachedItems)

	assert.Equal(t, 2, len(cachedItems))

	c.Clear()
	err = c.SetMultiply(map[string]any{}, fiveMinute)
	assert.Nil(t, err)

	c.Clear()
	oneNanosecond := 1 * time.Nanosecond
	err = c.SetMultiply(data, oneNanosecond)

	time.Sleep(1 * time.Second)

	assert.Nil(t, c.Get("example3", nil))
}

func TestFileSimple_InvalidJson(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFileSimple(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := cache.NewFileItem("example")
	_, err = item.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item.ExpiresAfter(5 * time.Minute)

	item2 := cache.NewFileItem("example2")
	_, err = item2.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item2.ExpiresAfter(5 * time.Minute)

	item3 := cache.NewFileItem("example3")
	_, err = item3.Set("Hello, world!", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item value")

	item3.ExpiresAfter(5 * time.Minute)

	err = c.Set(item2.GetKey(), item2.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Set(item3.GetKey(), item3.Get(), standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")

	err = c.DeleteMultiply("example", "example2")
	assert.NoError(t, err, "Failed to delete cache items")

	cachedItems := c.GetMultiply([]string{"example", "example2", "example3"}, nil)
	assert.NotEmpty(t, cachedItems)

	assert.Equal(t, 1, len(cachedItems))
}

func TestFileSimple_DeleteItem(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFileSimple(dir)
	assert.NoError(t, err, "Failed to create file cache")

	err = c.Set("item_to_delete", "Delete me", standards.KeepTTL)
	assert.NoError(t, err, "Failed to save cache item")

	err = c.Delete("item_to_delete")
	assert.NoError(t, err, "Failed to delete cache item")

	cachedItem := c.Get("item_to_delete", nil)
	assert.Nil(t, cachedItem)

	assert.False(t, c.Has("non-exists-key"))
}

func TestFileSimple_Save(t *testing.T) {
	t.Run("Save should handle json.Marshal error", func(t *testing.T) {
		tempDir := t.TempDir()
		fileCache, err := cache.NewFileSimple(tempDir)
		assert.NoError(t, err)

		err = fileCache.Set("invalid-json", make(chan int), standards.KeepTTL)
		assert.Error(t, err)
	})

	t.Run("should return error when os.ReadDir fails 2", func(t *testing.T) {
		dir := "/non-existent-directory"
		_, err := cache.NewFileSimple(dir)
		assert.Error(t, err)
	})

	t.Run("should return error when os.Remove fails", func(t *testing.T) {
		tempDir := t.TempDir()

		filePath := filepath.Join(tempDir, "test-file")
		err := os.WriteFile(filePath, []byte("data"), 0000)
		_ = os.Chmod(filePath, 0000)
		assert.NoError(t, err)

		fileCache, err := cache.NewFileSimple(tempDir)

		os.RemoveAll(tempDir)
		dir := filepath.Join("/non-exist", "cache_test")
		fileCache.Dir = dir
		err = fileCache.Clear()
		assert.Error(t, err)

		if err != nil {
			assert.Contains(t, err.Error(), "no such file or directory")
		}

		_ = os.Chmod(filePath, 0644)

	})

	t.Run("should return error when os.Remove fails (file not exists)", func(t *testing.T) {
		tempDir := t.TempDir()

		filePath := filepath.Join(tempDir, "test-file")
		err := os.WriteFile(filePath, []byte("data"), 0000) // Nastavení nulových oprávnění
		_ = os.Chmod(filePath, 0000)
		assert.NoError(t, err)

		fileCache, err := cache.NewFileSimple(tempDir)

		os.Remove(filePath)
		err = fileCache.Clear()

		if err != nil {
			assert.Contains(t, err.Error(), "no such file or directory")
		}

		_ = os.Chmod(filePath, 0644)
	})

	t.Run("Json errors", func(t *testing.T) {

		t.Run("Save should handle json.Marshal error", func(t *testing.T) {
			tempDir := t.TempDir()
			fileCache, err := cache.NewFileSimple(tempDir)
			assert.NoError(t, err)

			item := cache.NewFileItem("invalid-json")
			item.Set(make(chan int), standards.KeepTTL)

			oneMinute := 1 * time.Minute
			err = fileCache.SetMultiply(map[string]any{"invalid-json": item.Get()}, oneMinute)
			assert.Error(t, err)
		})
	})
}
