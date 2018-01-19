package sinks

import (
	"os"
)

// The BasicFileSink creates a file with the specified destinationFilename,
// overwriting if necessary and returns an io.Writer that is that file
type BasicFileSink struct {
	destinationFilename string
	file                *os.File
}

func NewBasicFileSink(destinationFilename string) *BasicFileSink {
	return &BasicFileSink{
		destinationFilename: destinationFilename,
	}
}

func (f *BasicFileSink) Open() error {
	file, err := os.Create(f.destinationFilename)
	f.file = file
	return err
}
func (f *BasicFileSink) Close() error {
	return f.file.Close()
}

func (f *BasicFileSink) Write(b []byte) (int, error) {
	return f.file.Write(b)
}
