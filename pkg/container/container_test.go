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
