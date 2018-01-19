package sinks

import "io"

type SinkParams struct {
	InputFilename string
}

type SinkFactory interface {
	Sink(params SinkParams) (Sink, error)
}

type Sink interface {
	Open() error
	io.Closer
	io.Writer
}
