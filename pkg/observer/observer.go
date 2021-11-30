package observer

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/utils"
	"github.com/reactivex/rxgo/v2"
)

func areWriteEvents(item interface{}) bool {
	event := item.(fsnotify.Event)
	return event.Op&fsnotify.Write == fsnotify.Write
}

func areChanges(item interface{}, hashesCh chan map[string]string) bool {
	event := item.(fsnotify.Event)
	hashes := <-hashesCh
	filename := event.Name
	newHash := utils.CreateFileHash(filename)
	hash, isInMap := hashes[filename]
	isDiff := true

	if isInMap {
		isDiff = hash != newHash
	}

	hashes[filename] = newHash
	hashesCh <- hashes
	return isDiff
}

func CreateObserverChannel() chan rxgo.Item {
	return make(chan rxgo.Item)
}

func ObserveItem(observableCh chan<- rxgo.Item, value interface{}) {
	observableCh <- rxgo.Item{V: value}
}

func AutoBuild(observableCh <-chan rxgo.Item) {
	const DEBOUNCE_DURATION = time.Second
	hashesCh := make(chan map[string]string)
	hashesCh <- make(map[string]string)
	observable := rxgo.FromChannel(observableCh).
		Filter(areWriteEvents).
		Debounce(rxgo.WithDuration(DEBOUNCE_DURATION)).
		Filter(func(item interface{}) bool {
			return areChanges(item, hashesCh)
		})
	defer close(hashesCh)

	for item := range observable.Observe() {
		if item.Error() {
			log.Fatal(item.E.Error())
		}
		log.Println("modified file:", item.V.(fsnotify.Event).Name)
	}
}
