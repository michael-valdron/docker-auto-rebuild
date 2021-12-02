package utils_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/michael-valdron/docker-auto-rebuild/pkg/utils"
	"github.com/muesli/cache2go"
	"github.com/stretchr/testify/assert"
)

func TestGetCacheName(t *testing.T) {
	assert.Equal(t, "fileHashes", utils.GetCacheName())
}

func TestGetOrCreateCache(t *testing.T) {
	cache := utils.GetOrCreateCache()
	cache.Add("Foo", 10*time.Second, "Bar")
	ch := make(chan string)
	defer close(ch)
	defer cache.Flush()

	go func() {
		cache := utils.GetOrCreateCache()
		value, err := cache.Value("Foo")

		if err != nil {
			ch <- "foo"
			return
		}

		ch <- value.Data().(string)
	}()

	assert.Equal(t, "Bar", <-ch)
}

func TestSetCacheValue(t *testing.T) {
	cache := utils.GetOrCreateCache()
	defer cache.Flush()

	utils.SetCacheValue("Foo", "Bar")
	value, err := cache.Value("Foo")

	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, "Bar", value.Data())
}

func TestGetCacheValue(t *testing.T) {
	cache := utils.GetOrCreateCache()
	defer cache.Flush()

	cache.Add("Foo", 10*time.Second, "Bar")
	value := utils.GetCacheValue("Foo")

	if value == nil {
		assert.Fail(t, "nil was returned instead of 'Bar'.")
	} else {
		assert.Equal(t, "Bar", *value)
	}
}

func TestInitFilesCache(t *testing.T) {
	wd, _ := os.Getwd()
	fileOne := filepath.Join(wd, "cache_test.go")
	fileTwo := filepath.Join(wd, "foobar")
	var cache *cache2go.CacheTable
	var value *cache2go.CacheItem
	var err error

	utils.InitFilesCache(wd)
	cache = utils.GetOrCreateCache()
	defer cache.Flush()

	value, err = cache.Value(fileOne)
	if err != nil {
		panic(err.Error())
	}
	assert.NotNil(t, value.Data())

	_, err = cache.Value(fileTwo)
	if err != nil {
		assert.True(t, true)
	} else {
		assert.Fail(t, "Did not throw expected error for non-existant file.")
	}
}

func TestFlushCache(t *testing.T) {
	cache := utils.GetOrCreateCache()
	var value *cache2go.CacheItem
	var err error

	cache.Add("Foo", 10*time.Second, "Bar")

	value, err = cache.Value("Foo")

	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, "Bar", value.Data())

	utils.FlushCache()

	_, err = cache.Value("Foo")

	if err != nil {
		assert.True(t, true)
	} else {
		assert.Fail(t, "Did not throw expected error due to flushed cache.")
	}
}
