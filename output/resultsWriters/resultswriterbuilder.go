package resultsWriters

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"os"
)

type ResultsWriterBuilder struct {
	outputFormat         string // "table" (default), "csv", "json"
	writeToMultipleFiles bool   // If true, write results to a files "next to" each input file, with the specified extension
	outputFilename       string // If not writing to multiple files and not "", write all results to this file
}

func NewResultsWriterBuilder(outputFormat string, writeToMultipleFiles bool, outputFilename string) *ResultsWriterBuilder {
	return &ResultsWriterBuilder{
		outputFormat:         outputFormat,
		writeToMultipleFiles: writeToMultipleFiles,
		outputFilename:       outputFilename,
	}
}

func (b *ResultsWriterBuilder) Build() (ResultsWriter, error) {
	var (
		formatter           formatters.ClassifyResultsFormatter
		outputFileExtension string
		resultsWriter       ResultsWriter
		writerFactory       sinks.SinkFactory
	)

	if b.writeToMultipleFiles && b.outputFilename != "" {
		return nil, errors.New("The --multiple-files (-m) and --output-file (-o) options cannot be used together.")
	}

	switch b.outputFormat {
	case "table":
		formatter = formatters.NewTableClassifyResultsFormatter()
		outputFileExtension = ".tab"
	case "csv":
		formatter = formatters.NewCSVClassifyResultsFormatter()
		outputFileExtension = ".csv"
	case "json":
		formatter = formatters.NewJsonClassifyResultsFormatter()
		outputFileExtension = ".json"
	default:
		// DocOpt doesn't do validation of these values for us, so we need to catch invalid values here
		return nil, errors.New(fmt.Sprintf("Unknown output format '%s'. Allowed values are: csv, table, json.", b.outputFormat))
	}

	if b.writeToMultipleFiles {
		// Write output to a file "next to" each input file, with the specified extension
		writerFactory = sinks.NewExtensionSwappingFileSinkFactory(outputFileExtension)
		resultsWriter = NewIndividualResultsWriter(writerFactory, formatter)
	} else if b.outputFilename != "" {
		// Write output to a single "combined results" file with the specified filename
		resultsWriter = NewCombinedResultsWriter(sinks.NewBasicFileSink(afero.NewOsFs(), b.outputFilename), formatter)
	} else {
		// Write output to the console
		resultsWriter = NewCombinedResultsWriter(sinks.NewBasicWriterSink(os.Stdout), formatter)
	}

	return resultsWriter, nil
}
