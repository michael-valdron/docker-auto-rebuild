package observer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/container"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/observer"
	"github.com/michael-valdron/docker-auto-rebuild/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateObserverChannel(t *testing.T) {
	ch := observer.CreateObserverChannel()
	defer close(ch)

	assert.True(t, true)
}

func TestObserveItem(t *testing.T) {
	ch := observer.CreateObserverChannel()

	go observer.ObserveItem(ch, "test")

	assert.Equal(t, "test", (<-ch).V.(string))
}

func TestAreWriteEvents(t *testing.T) {
	eventOne := fsnotify.Event{
		Name: "test.txt",
		Op:   fsnotify.Write,
	}
	eventTwo := fsnotify.Event{
		Name: "test.txt",
		Op:   fsnotify.Chmod,
	}

	assert.True(t, observer.AreWriteEvents(eventOne))
	assert.False(t, observer.AreWriteEvents(eventTwo))
}

func TestAreChanges(t *testing.T) {
	event := fsnotify.Event{
		Name: "test.txt",
		Op:   fsnotify.Write,
	}
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, event.Name)
	fs, err := os.Create(filePath)

	if err != nil {
		panic(fmt.Sprintf("%s cannot be created.", event.Name))
	}

	defer os.Remove(filePath)
	defer fs.Close()

	fs.WriteString("Hello World\n")

	utils.InitFilesCache(wd)

	assert.False(t, observer.AreChanges(event))

	fs.WriteString("Foobar\n")

	assert.True(t, observer.AreChanges(event))
}

func TestJustStringArray(t *testing.T) {
	compDirs := []string{"hello-world", "nodejs", "django"}
	observable := observer.JustStringArray(compDirs)
	result := []string{}

	for item := range observable.Observe() {
		if item.Error() {
			panic(item.E.Error())
		}

		result = append(result, item.V.(string))
	}

	assert.Equal(t, compDirs[0], result[0])
	assert.Equal(t, compDirs[1], result[1])
	assert.Equal(t, compDirs[2], result[2])
}

func TestComponentFile(t *testing.T) {
	compDir := "hello-world"
	event := fsnotify.Event{
		Name: "test.txt",
		Op:   fsnotify.Write,
	}
	compFile := observer.BuildComponentFile(compDir, event)

	assert.Equal(t, compDir, compFile.ComponentPath)
	assert.Equal(t, event.Name, compFile.Filename)
}

func TestComponentFiles(t *testing.T) {
	compDirs := []string{"hello-world", "nodejs", "django"}
	event := fsnotify.Event{
		Name: "test.txt",
		Op:   fsnotify.Write,
	}
	compFiles := []container.ComponentFile{}
	observable := observer.BuildComponentFiles(compDirs, event)

	for item := range observable.Observe() {
		if item.Error() {
			panic(item.E.Error())
		}

		compFile := item.V.(container.ComponentFile)
		compFiles = append(compFiles, compFile)
		assert.Equal(t, event.Name, compFile.Filename)
	}

	assert.Equal(t, compDirs[0], compFiles[0].ComponentPath)
	assert.Equal(t, compDirs[1], compFiles[1].ComponentPath)
	assert.Equal(t, compDirs[2], compFiles[2].ComponentPath)
}
