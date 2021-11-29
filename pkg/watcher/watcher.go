package watcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

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
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	startWatching(watcher, stop, observeEvent)

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	done <- true
}
