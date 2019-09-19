package assertions

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/waives/surf/fs"
	"os"
	"testing"
)

func FileExists(t *testing.T, name string) {
	exists, _ := fs.DirectoryOrFileExists(name)
	if !exists {
		assert.Fail(t, fmt.Sprintf("File %s does not exist", name))
	}
}

func FileDoesNotExist(t *testing.T, name string) {
	exists, _ := fs.DirectoryOrFileExists(name)
	if exists {
		assert.Fail(t, fmt.Sprintf("File %s exists when it should not", name))
	}
}

func DirectoryExists(t *testing.T, name string) {
	exists, _ := fs.DirectoryOrFileExists(name)
	if !exists {
		assert.Fail(t, fmt.Sprintf("Directory %s does not exist", name))
	}
}

func DirectoryDoesNotExist(t *testing.T, name string) {
	exists, _ := fs.DirectoryOrFileExists(name)
	if exists {
		assert.Fail(t, fmt.Sprintf("Directory %s exists when it should not", name))
	}
}

func DirectoryOrFileHasPermissions(t *testing.T, name string, perm os.FileMode) {
	info, err := os.Stat(name)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, perm, info.Mode().Perm())
}
