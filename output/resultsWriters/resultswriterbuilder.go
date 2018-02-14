package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
	"github.com/spf13/afero"
	"os"
)

var OutputFormatExtensions = map[formatters.OutputFormat]string{
	formatters.Table: ".tab",
	formatters.Json:  ".json",
	formatters.Csv:   ".csv",
}

type OutputType int

const (
	MultipleFiles OutputType = iota
	SingleFile
	Stdout
)

func NewResultsWriterFor(params *config.RunParams) (ResultsWriter, error) {

	formatter, err := formatters.NewResultsFormatterFor(params)

	if err != nil {
		return nil, err
	}

	if params.MultiFileOut {
		// Write output to a file "next to" each input file, with the specified extension
		return NewResultsWriter(MultipleFiles, "", formatter)
	} else if params.OutputFile != "" {
		// Write output to a single "combined results" file with the specified filename
		return NewResultsWriter(SingleFile, params.OutputFile, formatter)
	} else {
		// Write output to the console
		return NewResultsWriter(Stdout, "", formatter)
	}
}

func NewResultsWriter(outputType OutputType, outputFilename string, formatter formatters.ResultsFormatter) (ResultsWriter, error) {
	var (
		outputFileExtension string
		resultsWriter       ResultsWriter
		writerFactory       sinks.SinkFactory
	)

	outputFileExtension = OutputFormatExtensions[formatter.Format()]

	switch outputType {
	case MultipleFiles:
		// Write output to a file "next to" each input file, with the specified extension
		writerFactory = sinks.NewExtensionSwappingFileSinkFactory(outputFileExtension)
		resultsWriter = NewIndividualResultsWriter(writerFactory, formatter)
	case SingleFile:
		// Write output to a single "combined results" file with the specified filename
		resultsWriter = NewCombinedResultsWriter(sinks.NewBasicFileSink(afero.NewOsFs(), outputFilename), formatter)
	case Stdout:
		// Write output to the console
		resultsWriter = NewCombinedResultsWriter(sinks.NewBasicWriterSink(os.Stdout), formatter)
	}
	return resultsWriter, nil
}
