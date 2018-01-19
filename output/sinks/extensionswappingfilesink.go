package sinks

import (
	"os"
	"path/filepath"
	"strings"
)

// The ExtensionSwappingFileSink creates a file adjacent to the specified inputFilename (but with the specified extension),
// overwriting if necessary and returns an io.Writer that is that file
type ExtensionSwappingFileSink struct {
	destinationFilename string
	file                *os.File
}

func newExtensionSwappingFileSink(fileExtension string, inputFilename string) *ExtensionSwappingFileSink {
	return &ExtensionSwappingFileSink{
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
	file, err := os.Create(f.destinationFilename)
	f.file = file
	return err
}

func (f *ExtensionSwappingFileSink) Close() error {
	return f.file.Close()
}

func (f *ExtensionSwappingFileSink) Write(b []byte) (int, error) {
	//TODO Call Open if not already called, or throw Error
	return f.file.Write(b)
}
