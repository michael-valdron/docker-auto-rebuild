package utils

import (
	"crypto/sha256"
	"io/ioutil"
	"log"
)

func ReadFileContents(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func CreateFileHash(filename string) string {
	content := ReadFileContents(filename)
	hasher := sha256.New()

	hasher.Write([]byte(filename))
	hasher.Write(content)

	return string(hasher.Sum(nil))
}
