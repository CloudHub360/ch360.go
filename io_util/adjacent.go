package io_util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type AdjacentFileProvider struct {
	Extension string
}

func (afp *AdjacentFileProvider) Provide(fullPath string) (io.WriteCloser, error) {
	filename := filepath.Base(fullPath)
	outFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + afp.Extension
	outFullPath := filepath.Join(filepath.Dir(fullPath), outFilename)
	return os.Create(outFullPath)
}
