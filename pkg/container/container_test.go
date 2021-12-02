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
	assert.True(t, container.IsFileInComponent("../../examples/docker-project/hello-world", "../../examples/docker-project/hello-world/Dockerfile"))
	assert.True(t, container.IsFileInComponent("../../examples/docker-project/nodejs", "../../examples/docker-project/nodejs/test/test.js"))
	assert.False(t, container.IsFileInComponent("../../examples/docker-project/hello-world", "../../examples/docker-project/nodejs/test/test.js"))
	assert.False(t, container.IsFileInComponent("../../examples/docker-project/nodejs", "../../examples/docker-project/hello-world/Dockerfile"))
}
