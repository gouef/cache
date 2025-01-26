package tests

import (
	"github.com/gouef/cache"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func setupCacheDir(t *testing.T) string {
	dir := filepath.Join(os.TempDir(), "cache_test")
	err := os.MkdirAll(dir, 0755)
	assert.NoError(t, err, "Failed to create test cache directory")
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return dir
}

func TestFile_SaveAndGetItem(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := cache.NewFileItem("example")
	_, err = item.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item.ExpiresAfter(5 * time.Minute)

	err = c.Save(item)
	assert.NoError(t, err, "Failed to save cache item")

	cachedItem := c.GetItem("example")
	assert.NotNil(t, cachedItem)

	if cachedItem != nil {
		assert.True(t, cachedItem.IsHit(), "Expected cache item to be a hit")
		assert.Equal(t, "Hello, world!", cachedItem.Get(), "Cache item value mismatch")
		assert.Equal(t, "example", cachedItem.GetKey())
	}

	// Invalid

	assert.Error(t, c.Save(nil))
}

func TestFile_SaveAndGetItems(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := cache.NewFileItem("example")
	_, err = item.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item.ExpiresAfter(5 * time.Minute)

	item2 := cache.NewFileItem("example2")
	_, err = item2.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item2.ExpiresAfter(5 * time.Minute)

	item3 := cache.NewFileItem("example3")
	_, err = item3.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item3.ExpiresAfter(1 * time.Second)

	err = c.SaveDeferred(item)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Save(item2)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Save(item3)
	assert.NoError(t, err, "Failed to save cache item")

	time.Sleep(1 * time.Second)

	assert.Nil(t, c.Commit())

	cachedItems := c.GetItems("example", "example2")
	assert.NotEmpty(t, cachedItems)

	assert.Equal(t, 2, len(cachedItems))

	// Invalid json
	mu := sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()

	data := []byte("{{,#)")

	os.WriteFile(filepath.Join(dir, item.GetKey()+".cache"), data, 0644)

	invalidItem := c.GetItem(item.GetKey())

	assert.Nil(t, invalidItem)

}

func TestFile_InvalidJson(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := cache.NewFileItem("example")
	_, err = item.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item.ExpiresAfter(5 * time.Minute)

	item2 := cache.NewFileItem("example2")
	_, err = item2.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item2.ExpiresAfter(5 * time.Minute)

	item3 := cache.NewFileItem("example3")
	_, err = item3.Set("Hello, world!")
	assert.NoError(t, err, "Failed to save cache item value")

	item3.ExpiresAfter(5 * time.Minute)

	err = c.SaveDeferred(item)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Save(item2)
	assert.NoError(t, err, "Failed to save cache item")
	err = c.Save(item3)
	assert.NoError(t, err, "Failed to save cache item")

	assert.Nil(t, c.Commit())

	err = c.DeleteItems("example", "example2")
	assert.NoError(t, err, "Failed to delete cache items")

	cachedItems := c.GetItems("example", "example2", "example3")
	assert.NotEmpty(t, cachedItems)

	assert.Equal(t, 1, len(cachedItems))
}

func TestFile_ItemExpiration(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := &cache.FileItem{
		Key:   "expiring_item",
		Value: "This will expire",
	}
	item.ExpiresAfter(1 * time.Second)

	err = c.Save(item)
	assert.NoError(t, err, "Failed to save cache item")

	cachedItem := c.GetItem("expiring_item")
	assert.NotNil(t, cachedItem)

	if cachedItem != nil {
		assert.True(t, cachedItem.IsHit(), "Expected cache item to be a hit before expiration")
		assert.Equal(t, "This will expire", cachedItem.Get(), "Cache item value mismatch before expiration")
	}

	time.Sleep(3 * time.Second)

	assert.Nil(t, item.Get())
	cachedItem = c.GetItem("expiring_item")
	assert.Nil(t, cachedItem, "Expected cache item to be a miss after expiration")
}

func TestFile_Clear(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	items := []string{"item1", "item2", "item3"}
	for _, key := range items {
		item := &cache.FileItem{
			Key:   key,
			Value: key + "_value",
		}
		err := c.Save(item)
		assert.NoError(t, err, "Failed to save cache item")
	}

	for _, key := range items {
		cachedItem := c.GetItem(key)
		assert.True(t, cachedItem.IsHit(), "Expected cache item to be a hit")
	}

	err = c.Clear()
	assert.NoError(t, err, "Failed to clear cache")

	for _, key := range items {
		cachedItem := c.GetItem(key)
		assert.Nil(t, cachedItem)
	}
}

func TestFile_DeleteItem(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := &cache.FileItem{
		Key:   "item_to_delete",
		Value: "Delete me",
	}
	err = c.Save(item)
	assert.NoError(t, err, "Failed to save cache item")

	err = c.DeleteItem("item_to_delete")
	assert.NoError(t, err, "Failed to delete cache item")

	cachedItem := c.GetItem("item_to_delete")
	assert.Nil(t, cachedItem)
}

func TestFile_HasItem(t *testing.T) {
	dir := setupCacheDir(t)
	c, err := cache.NewFile(dir)
	assert.NoError(t, err, "Failed to create file cache")

	item := &cache.FileItem{
		Key:   "existing_item",
		Value: "I exist",
	}
	err = c.Save(item)
	assert.NoError(t, err, "Failed to save cache item")

	assert.True(t, c.HasItem("existing_item"), "Expected cache to contain the item")

	assert.False(t, c.HasItem("non_existing_item"), "Expected cache to not contain the item")
}

func TestFile_NewFile(t *testing.T) {
	dir := filepath.Join("/non-exist", "cache_test")

	_, err := cache.NewFile(dir)
	assert.Error(t, err, "Failed to create file cache")
}

func TestFile_Save(t *testing.T) {
	t.Run("Save should handle json.Marshal error", func(t *testing.T) {
		tempDir := t.TempDir()
		fileCache, err := cache.NewFile(tempDir)
		assert.NoError(t, err)

		item := cache.NewFileItem("invalid-json")
		item.Set(make(chan int))

		err = fileCache.Save(item)
		assert.Error(t, err)
	})

	t.Run("should return error when os.ReadDir fails 2", func(t *testing.T) {
		dir := "/non-existent-directory"
		fileCache := &cache.File{Dir: dir}

		err := fileCache.Clear()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("should return error when os.Remove fails", func(t *testing.T) {
		tempDir := t.TempDir()

		filePath := filepath.Join(tempDir, "test-file")
		err := os.WriteFile(filePath, []byte("data"), 0000)
		_ = os.Chmod(filePath, 0000)
		assert.NoError(t, err)

		fileCache := &cache.File{Dir: tempDir}

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
		err := os.WriteFile(filePath, []byte("data"), 0000)
		_ = os.Chmod(filePath, 0000)
		assert.NoError(t, err)

		fileCache := &cache.File{Dir: tempDir}

		os.Remove(filePath)
		err = fileCache.Clear()

		if err != nil {
			assert.Contains(t, err.Error(), "no such file or directory")
		}

		_ = os.Chmod(filePath, 0644)
	})
}
