package commands

import (
	"github.com/CloudHub360/ch360.go/io_util"
	"io"
)

//go:generate mockery -name "WriterProvider|WriteCloserProvider"

type WriterProvider interface {
	Provide(fullPath string) (io.Writer, error)
}
type WriteCloserProvider interface {
	Provide(fullPath string) (io.WriteCloser, error)
}

// BasicWriterProvider just returns the provided io.Writer in any call to Provide
type BasicWriterProvider struct {
	dest io.Writer
}

func NewBasicWriterProvider(dest io.Writer) *BasicWriterProvider {
	return &BasicWriterProvider{
		dest: dest,
	}
}

func (f *BasicWriterProvider) Provide(fullPath string) (io.Writer, error) {
	return f.dest, nil
}

// NewAutoClosingWriterProvider wraps any io.WriteClosers returned by its underlying WriterProvider in an
// io_util.AutoCloser.
type AutoClosingWriterProvider struct {
	underlying WriteCloserProvider
}

func (f *AutoClosingWriterProvider) Provide(fullPath string) (io.Writer, error) {
	if fullPath == "" {
		return nil, nil
	}

	outWriter, err := f.underlying.Provide(fullPath)

	if err != nil {
		return nil, err
	}

	return &io_util.AutoCloser{
		Underlying: outWriter,
	}, nil
}

func NewAutoClosingWriterProvider(underlying WriteCloserProvider) *AutoClosingWriterProvider {
	return &AutoClosingWriterProvider{
		underlying: underlying,
	}
}
