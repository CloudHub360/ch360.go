package sinks

import (
	"github.com/spf13/afero"
)

// The BasicFileSink creates a file with the specified destinationFilename,
// overwriting if necessary and returns an io.Writer that is that file
type BasicFileSink struct {
	fileSystem          afero.Fs
	destinationFilename string
	file                afero.File
}

func NewBasicFileSink(fileSystem afero.Fs, destinationFilename string) *BasicFileSink {
	return &BasicFileSink{
		fileSystem:          fileSystem,
		destinationFilename: destinationFilename,
	}
}

func (f *BasicFileSink) Open() error {
	file, err := f.fileSystem.Create(f.destinationFilename)
	f.file = file
	return err
}
func (f *BasicFileSink) Close() error {
	return f.file.Close()
}

func (f *BasicFileSink) Write(b []byte) (int, error) {
	return f.file.Write(b)
}
