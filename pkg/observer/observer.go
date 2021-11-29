package observer

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/reactivex/rxgo/v2"
)

func CreateObserverChannel() chan rxgo.Item {
	return make(chan rxgo.Item)
}

func ObserveItem(observableCh chan<- rxgo.Item, value interface{}) {
	observableCh <- rxgo.Item{V: value}
}

func AutoBuild(observableCh <-chan rxgo.Item) {
	observable := rxgo.FromChannel(observableCh).Filter(func(item interface{}) bool {
		event := item.(fsnotify.Event)
		return event.Op&fsnotify.Write == fsnotify.Write
	})

	for item := range observable.Observe() {
		if item.Error() {
			log.Fatal(item.E.Error())
		}
		log.Println("modified file:", item.V.(fsnotify.Event).Name)
	}
}
