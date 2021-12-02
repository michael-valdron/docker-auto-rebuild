package observer

import (
	"context"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/container"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/utils"
	"github.com/reactivex/rxgo/v2"
)

func CreateObserverChannel() chan rxgo.Item {
	return make(chan rxgo.Item)
}

func ObserveItem(observableCh chan<- rxgo.Item, value interface{}) {
	observableCh <- rxgo.Item{V: value}
}

func AreWriteEvents(item interface{}) bool {
	event := item.(fsnotify.Event)
	return event.Op&fsnotify.Write == fsnotify.Write
}

func AreChanges(item interface{}) bool {
	event := item.(fsnotify.Event)
	filename := event.Name
	newHash := utils.CreateFileHash(filename)
	hash := utils.GetOrCreateFileHash(filename)
	isDiff := hash != newHash

	if isDiff {
		utils.SetCacheValue(filename, newHash)
	}

	return isDiff
}

func JustStringArray(arr []string) rxgo.Observable {
	items := make([]interface{}, len(arr))
	for i := 0; i < len(arr); i++ {
		items[i] = arr[i]
	}
	return rxgo.Just(items...)()
}

func BuildComponentFile(compDirItem interface{}, eventItem interface{}) container.ComponentFile {
	event := eventItem.(fsnotify.Event)
	compDir := compDirItem.(string)
	return container.ComponentFile{
		ComponentPath: compDir,
		Filename:      event.Name,
	}
}

func BuildComponentFiles(compDirs []string, eventItem interface{}) rxgo.Observable {
	return JustStringArray(compDirs).
		Map(func(c context.Context, compDirItem interface{}) (interface{}, error) {
			return BuildComponentFile(compDirItem, eventItem), c.Err()
		})
}

func AutoBuild(observableCh <-chan rxgo.Item, workingDir string) {
	const DEBOUNCE_DURATION = 3 * time.Second
	compDirs := container.GetComponentDirectories(workingDir)
	fileEvents := rxgo.FromChannel(observableCh).
		Filter(AreWriteEvents).
		Debounce(rxgo.WithDuration(DEBOUNCE_DURATION)).
		Filter(AreChanges)
	buildDirs := fileEvents.
		FlatMap(func(item rxgo.Item) rxgo.Observable {
			return BuildComponentFiles(compDirs, item.V)
		}).
		Filter(func(item interface{}) bool {
			return item.(container.ComponentFile).IsFileInComponent()
		}).
		Map(func(c context.Context, item interface{}) (interface{}, error) {
			return item.(container.ComponentFile).ComponentPath, c.Err()
		})

	for buildDir := range buildDirs.Observe() {
		if buildDir.Error() {
			log.Fatal(buildDir.E.Error())
		}
		container.RunBuild(buildDir.V.(string))
	}
}
