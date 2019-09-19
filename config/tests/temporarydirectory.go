package tests

import (
	"fmt"
	"github.com/waives/surf/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type TemporaryDirectory struct {
	guid string
	path string
}

func NewTemporaryDirectory() *TemporaryDirectory {
	return &TemporaryDirectory{
		guid: fmt.Sprintf("%v", time.Now().UTC().UnixNano()),
	}
}

func (dir *TemporaryDirectory) Path() string {
	if dir.path == "" {
		tmpDir, _ := ioutil.TempDir("", "fakehome")
		dir.path = filepath.Join(tmpDir, dir.guid)
	}
	return dir.path
}

func (dir *TemporaryDirectory) Create() {
	os.MkdirAll(dir.Path(), config.DirRWPermissions)
}

func (dir *TemporaryDirectory) Destroy() {
	os.RemoveAll(dir.Path())
}
