package fakes

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FakeHomeDirectoryPathGetter struct {
	guid string
	path string
}

func NewFakeHomeDirectoryPathGetter() *FakeHomeDirectoryPathGetter {
	return &FakeHomeDirectoryPathGetter{
		guid: fmt.Sprintf("%v", time.Now().UTC().UnixNano()),
	}
}

func (dir *FakeHomeDirectoryPathGetter) Path() string {
	if dir.path == "" {
		tmpDir, _ := ioutil.TempDir("", "fakehome")
		dir.path = filepath.Join(tmpDir, dir.guid)
	}
	return dir.path
}

func (dir *FakeHomeDirectoryPathGetter) Create() {
	os.MkdirAll(dir.Path(), config.DirRWPermissions)
}

func (dir *FakeHomeDirectoryPathGetter) Destroy() {
	os.RemoveAll(dir.Path())
}
