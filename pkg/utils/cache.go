package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/muesli/cache2go"
)

func GetCacheName() string {
	return "fileHashes"
}

func GetOrCreateCache() *cache2go.CacheTable {
	return cache2go.Cache(GetCacheName())
}

func SetCacheValue(key string, value string) {
	cache := GetOrCreateCache()
	cache.Add(key, 10*time.Minute, value)
}

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
