package commands

import (
	"github.com/CloudHub360/ch360.go/io_util"
	"io"
)

type WriteCloserProvider func(fullPath string) (io.WriteCloser, error)

var DummyWriteCloserProvider = func(dest io.WriteCloser) WriteCloserProvider {
	return func(fullPath string) (io.WriteCloser, error) {
		return dest, nil
	}
}

var AutoClosingWriteCloserProvider = func(underlying WriteCloserProvider) WriteCloserProvider {
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
