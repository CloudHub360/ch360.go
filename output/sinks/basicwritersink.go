package sinks

import (
	"github.com/CloudHub360/ch360.go/fs"
	"io"
)

type BasicWriterSink struct {
	writer io.Writer
}

func NewBasicWriterSink(writer io.Writer) *BasicWriterSink {
	return &BasicWriterSink{
		writer: writer,
	}
}

func (f *BasicWriterSink) Open() error {
	return nil
}

func (f *BasicWriterSink) Close() error {
	// the underlying writer could well be a file, in which case we should
	// try to close it here.
	fs.TryClose(f.writer)

	return nil
}

func (f *BasicWriterSink) Write(b []byte) (int, error) {
	return f.writer.Write(b)
}
