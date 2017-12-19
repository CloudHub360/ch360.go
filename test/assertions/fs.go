package assertions

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func FileExists(t *testing.T, name string) {
	exists, _ := directoryOrFileExists(name)
	if !exists {
		assert.Fail(t, fmt.Sprintf("File %s does not exist", name))
	}
}

func FileDoesNotExist(t *testing.T, name string) {
	exists, _ := directoryOrFileExists(name)
	if exists {
		assert.Fail(t, fmt.Sprintf("File %s exists when it should not", name))
	}
}

func DirectoryExists(t *testing.T, name string) {
	exists, _ := directoryOrFileExists(name)
	if !exists {
		assert.Fail(t, fmt.Sprintf("Directory %s does not exist", name))
	}
}

func DirectoryDoesNotExist(t *testing.T, name string) {
	exists, _ := directoryOrFileExists(name)
	if exists {
		assert.Fail(t, fmt.Sprintf("Directory %s exists when it should not", name))
	}
}

func directoryOrFileExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// directory does not exist
			return false, nil
		} else {
			// other error
			return false, err
		}
	}

	return true, nil
}
