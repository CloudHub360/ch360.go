package sinks

import "github.com/spf13/afero"

type ExtensionSwappingFileWriterFactory struct {
	fileExtension string
}

// The ExtensionSwappingFileWriterFactory returns a new ExtensionSwappingFileSink (pointing to a new destination file) each time Sink is called
func NewExtensionSwappingFileSinkFactory(fileExtension string) *ExtensionSwappingFileWriterFactory {
	return &ExtensionSwappingFileWriterFactory{
		fileExtension: fileExtension,
	}
}

func (p *ExtensionSwappingFileWriterFactory) Sink(params SinkParams) (Sink, error) {
	return NewExtensionSwappingFileSink(afero.NewOsFs(), p.fileExtension, params.InputFilename), nil
}
