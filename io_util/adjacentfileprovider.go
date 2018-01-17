package io_util

import (
	"os"
	"path/filepath"
	"strings"
)

// AdjacentFileProvider creates a file in the same directory as the fullPath specified in Provide
// but with the extension as given in the type's field.
type AdjacentFileProvider struct {
	Extension string
}

func (afp *AdjacentFileProvider) Provide(fullPath string) (*os.File, error) {
	// the filename without the path
	filename := filepath.Base(fullPath)

	// the filename with extension replaced
	outFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + afp.Extension

	// the output full path
	outFullPath := filepath.Join(filepath.Dir(fullPath), outFilename)
	return os.Create(outFullPath)
}
