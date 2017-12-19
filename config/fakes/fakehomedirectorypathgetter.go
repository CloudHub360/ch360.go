package fakes

import (
	"github.com/CloudHub360/ch360.go/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FakeHomeDirectoryPathGetter struct {
	Guid string
	path string
}

func (dir *FakeHomeDirectoryPathGetter) GetPath() string {
	if dir.path == "" {
		tmpDir, _ := ioutil.TempDir("", "fakehome")
		dir.path = filepath.Join(tmpDir, dir.Guid)
	}
	return dir.path
}

func (dir *FakeHomeDirectoryPathGetter) Create() {
	os.MkdirAll(dir.GetPath(), config.DirRWPermissions)
}

func (dir *FakeHomeDirectoryPathGetter) Destroy() {
	os.RemoveAll(dir.GetPath())
}
