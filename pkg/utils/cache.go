package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/muesli/cache2go"
)

// Gets the name identity for the cache.
func GetCacheName() string {
	return "fileHashes"
}

// Creates the cache if it does not exist else gets the active cache.
func GetOrCreateCache() *cache2go.CacheTable {
	return cache2go.Cache(GetCacheName())
}

// Set a value to a given key in the cache.
func SetCacheValue(key string, value string) {
	cache := GetOrCreateCache()
	cache.Add(key, 10*time.Minute, value)
}

// Gets a value from the cache by a given key.
func GetCacheValue(key string) *string {
	cache := GetOrCreateCache()
	var value *string = nil

	if cache.Exists(key) {
		itemPtr, err := cache.Value(key)
		if err != nil {
			log.Fatal(err)
		} else {
			item := itemPtr.Data().(string)
			value = &item
		}
	}

	return value
}

// Initializes the cache with files with their paths as the keys and their hashes as the values.
//
// See: utils.GetOrCreateFileHash
func InitFilesCache(workingDir string) {
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			GetOrCreateFileHash(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func FlushCache() {
	cache := GetOrCreateCache()
	cache.Flush()
}
