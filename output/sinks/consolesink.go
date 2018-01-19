package sinks

import (
	"os"
)

type ConsoleSink struct{}

func (f *ConsoleSink) Open() error {
	return nil
}

func (f *ConsoleSink) Close() error {
	return nil
}

func (f *ConsoleSink) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}
