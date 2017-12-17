package fakes

import (
	"os"
	"path/filepath"
)

type FakeHomeDirectoryPathGetter struct {
	Guid string
}

func (provider *FakeHomeDirectoryPathGetter) GetPath() string {
	path := filepath.Join(
		os.Getenv("GOPATH"),
		"src", "github.com", "CloudHub360", "ch360.go", "test", "output", provider.Guid)
	return path
}
