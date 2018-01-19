package sinks

import (
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
)

// The ExtensionSwappingFileSink creates a file adjacent to the specified inputFilename (but with the specified extension),
// overwriting if necessary and returns an io.Writer that is that file
type ExtensionSwappingFileSink struct {
	fileSystem          afero.Fs
	destinationFilename string
	file                afero.File
}

func NewExtensionSwappingFileSink(fileSystem afero.Fs, fileExtension string, inputFilename string) *ExtensionSwappingFileSink {
	return &ExtensionSwappingFileSink{
		fileSystem:          fileSystem,
		destinationFilename: replaceFileExtension(inputFilename, fileExtension),
	}
}

func replaceFileExtension(fullPath string, newFileExtension string) string {
	// the destinationFilename without the path
	filename := filepath.Base(fullPath)

	// the destinationFilename with extension replaced
	outFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + newFileExtension

	// the output full path
	outFullPath := filepath.Join(filepath.Dir(fullPath), outFilename)
	return outFullPath
}

func (f *ExtensionSwappingFileSink) Open() error {
	file, err := f.fileSystem.Create(f.destinationFilename)
	f.file = file
	return err
}

func (f *ExtensionSwappingFileSink) Close() error {
	return f.file.Close()
}

func (f *ExtensionSwappingFileSink) Write(b []byte) (int, error) {
	return f.file.Write(b)
}
