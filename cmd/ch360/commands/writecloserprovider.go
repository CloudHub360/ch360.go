package commands

import (
	"github.com/CloudHub360/ch360.go/io_util"
	"io"
)

type WriteCloserProvider func(fullPath string) (io.WriteCloser, error)

// DummyWriteCloserProvider just returns the provided io.WriteCloser in any call to the
func NewDummyWriteCloserProvider(dest io.WriteCloser) WriteCloserProvider {
	return func(fullPath string) (io.WriteCloser, error) {
		return dest, nil
	}
}

// NewAutoClosingWriteCloserProvider wraps any io.WriteClosers returned by its underlying WriteCloserProvider in an
// io_util.AutoCloser.
func NewAutoClosingWriteCloserProvider(underlying WriteCloserProvider) WriteCloserProvider {
	return func(fullPath string) (io.WriteCloser, error) {
		if fullPath == "" {
			return nil, nil
		}

		outWriter, err := underlying(fullPath)

		if err != nil {
			return nil, err
		}

		return &io_util.AutoCloser{
			Underlying: outWriter,
		}, nil
	}
}
