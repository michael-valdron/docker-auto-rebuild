package watcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func startWatching(watcher *fsnotify.Watcher, stop <-chan bool) {
	go func() {
		for {
			select {
			case <-stop:
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}

func Watch(stop <-chan bool, done chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	startWatching(watcher, stop)

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	done <- true
}
