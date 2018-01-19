package sinks

import "io"

type SinkParams struct {
	InputFilename string
}

//go:generate mockery -name SinkFactory
type SinkFactory interface {
	Sink(params SinkParams) (Sink, error)
}

//go:generate mockery -name Sink
type Sink interface {
	Open() error
	io.Closer
	io.Writer
}
