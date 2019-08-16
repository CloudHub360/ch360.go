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
	fs.TryClose(f.writer)

	return nil
}

func (f *BasicWriterSink) Write(b []byte) (int, error) {
	return f.writer.Write(b)
}
