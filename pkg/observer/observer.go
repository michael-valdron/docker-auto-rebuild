package observer

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/container"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/utils"
	"github.com/reactivex/rxgo/v2"
)

func areWriteEvents(item interface{}) bool {
	event := item.(fsnotify.Event)
	return event.Op&fsnotify.Write == fsnotify.Write
}

func areChanges(item interface{}, cache map[string]string) bool {
	event := item.(fsnotify.Event)
	filename := event.Name
	newHash := utils.CreateFileHash(filename)
	hash, isInMap := cache[filename]
	isDiff := true

	if isInMap {
		isDiff = hash != newHash
	}

	cache[filename] = newHash

	return isDiff
}

func CreateObserverChannel() chan rxgo.Item {
	return make(chan rxgo.Item)
}

func ObserveItem(observableCh chan<- rxgo.Item, value interface{}) {
	observableCh <- rxgo.Item{V: value}
}

func AutoBuild(observableCh <-chan rxgo.Item, workingDir string) {
	const DEBOUNCE_DURATION = time.Second
	cache := utils.CreateFilesCache(workingDir)
	fileEvents := rxgo.FromChannel(observableCh).
		Filter(areWriteEvents).
		Debounce(rxgo.WithDuration(DEBOUNCE_DURATION))

	for item := range fileEvents.Observe() {
		if item.Error() {
			log.Fatal(item.E.Error())
		}
		if areChanges(item.V, cache) {
			container.RunBuild()
		}
	}
}
