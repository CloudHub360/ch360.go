package fakes

import (
	"os"
	"path/filepath"
)

type FakeHomeDirectoryPathGetter struct {
	Guid string
}

func (dir *FakeHomeDirectoryPathGetter) GetPath() string {
	path := filepath.Join(
		os.Getenv("GOPATH"),
		"src", "github.com", "CloudHub360", "ch360.go", "test", "output", dir.Guid)
	return path
}

func (dir *FakeHomeDirectoryPathGetter) Create() {
	os.MkdirAll(dir.GetPath(), 600)
}

func (dir *FakeHomeDirectoryPathGetter) Destroy() {
	os.RemoveAll(dir.GetPath())
}
