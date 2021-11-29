package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func addRecursiveWatch(watcher *fsnotify.Watcher, rootPath string) error {
	err := watcher.Add(rootPath)
	if err != nil {
		return err
	}

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Mode().IsDir() {
			err = watcher.Add(path)
		}

		return err
	})
}

func startWatching(watcher *fsnotify.Watcher, stop <-chan bool, observeEvent func(interface{})) {
	go func() {
		for {
			select {
			case <-stop:
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				observeEvent(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}

func Watch(stop <-chan bool, done chan<- bool, observeEvent func(interface{})) {
	watcher, err := fsnotify.NewWatcher()
	var workerDir string
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	startWatching(watcher, stop, observeEvent)

	workerDir, _ = os.Getwd()
	err = addRecursiveWatch(watcher, workerDir)
	if err != nil {
		log.Fatal(err)
	}
	done <- true
}
