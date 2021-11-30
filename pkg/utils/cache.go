package utils

import (
	"log"
	"os"
	"path/filepath"
)

func CreateFilesCache(workingDir string) map[string]string {
	cache := make(map[string]string)
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			cache[path] = CreateFileHash(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return cache
}
