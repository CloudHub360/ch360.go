package sinks

import (
	"io"
)

type ConsoleSink struct {
	writer io.Writer
}

func NewConsoleSink(writer io.Writer) *ConsoleSink {
	return &ConsoleSink{
		writer: writer,
	}
}

func (f *ConsoleSink) Open() error {
	return nil
}

func (f *ConsoleSink) Close() error {
	return nil
}

func (f *ConsoleSink) Write(b []byte) (int, error) {
	return f.writer.Write(b)
}
