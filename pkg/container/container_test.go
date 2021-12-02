package container_test

import (
	"testing"

	"github.com/michael-valdron/docker-auto-rebuild/pkg/container"
	"github.com/stretchr/testify/assert"
)

func TestIsComponentDirectory(t *testing.T) {
	assert.True(t, container.IsComponentDirectory("../../examples/docker-project/hello-world"))
	assert.True(t, container.IsComponentDirectory("../../examples/docker-project/nodejs"))
	assert.False(t, container.IsComponentDirectory("../../examples/docker-project"))
}

func TestGetComponentDirectories(t *testing.T) {
	assert.Equal(t, 2, len(container.GetComponentDirectories("../../examples/docker-project")))
}

func TestIsFileInComponent(t *testing.T) {
	compFileOne := container.ComponentFile{
		ComponentPath: "../../examples/docker-project/hello-world",
		Filename:      "../../examples/docker-project/hello-world/Dockerfile",
	}
	compFileTwo := container.ComponentFile{
		ComponentPath: "../../examples/docker-project/nodejs",
		Filename:      "../../examples/docker-project/nodejs/test/test.js",
	}
	compFileThree := container.ComponentFile{
		ComponentPath: "../../examples/docker-project/hello-world",
		Filename:      "../../examples/docker-project/nodejs/test/test.js",
	}
	compFileFour := container.ComponentFile{
		ComponentPath: "../../examples/docker-project/nodejs",
		Filename:      "../../examples/docker-project/hello-world/Dockerfile",
	}

	assert.True(t, compFileOne.IsFileInComponent())
	assert.True(t, compFileTwo.IsFileInComponent())
	assert.False(t, compFileThree.IsFileInComponent())
	assert.False(t, compFileFour.IsFileInComponent())
}
