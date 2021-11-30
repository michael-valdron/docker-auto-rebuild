package utils_test

import (
	"path"
	"testing"

	"github.com/michael-valdron/docker-auto-rebuild/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestReadFileContents(t *testing.T) {
	content := string(utils.ReadFileContents("../../examples/docker-project/hello-world/Dockerfile"))
	assert.Equal(t, "FROM hello-world\n", content)
}

func TestCreateFileHash(t *testing.T) {
	projectDir := "../.."
	dirOne := path.Join(projectDir, "examples/docker-project/hello-world/Dockerfile")
	dirTwo := path.Join(projectDir, "examples/docker-project/nodejs/Dockerfile")
	hashOne := utils.CreateFileHash(dirOne)
	hashTwo := utils.CreateFileHash(dirTwo)
	assert.NotEqual(t, hashOne, hashTwo)
}
