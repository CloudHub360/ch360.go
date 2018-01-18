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

// BasicWriterFactory just returns the provided io.Writer in any call to Provide
type BasicWriterFactory struct {
	dest io.Writer
}

func NewBasicWriterFactory(dest io.Writer) *BasicWriterFactory {
	return &BasicWriterFactory{
		dest: dest,
	}
}

func (f *BasicWriterFactory) Provide(fullPath string) (io.Writer, error) {
	return f.dest, nil
}

// NewAutoClosingWriterFactory wraps any io.WriteClosers returned by its underlying WriterProvider in an
// io_util.AutoCloser.
type AutoClosingWriterFactory struct {
	underlying WriteCloserProvider
}

func (f *AutoClosingWriterFactory) Provide(fullPath string) (io.Writer, error) {
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

func NewAutoClosingWriterFactory(underlying WriteCloserProvider) *AutoClosingWriterFactory {
	return &AutoClosingWriterFactory{
		underlying: underlying,
	}
}
