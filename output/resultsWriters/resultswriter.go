package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/fs"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
)

//go:generate mockery -name ResultsWriter
type ResultsWriter interface {
	Start() error
	WriteResult(filename string, result interface{}) error
	Finish() error
}

// NewExtractionResultsWriter constructs a ResultsWriter configured for extraction.
func NewExtractionResultsWriter(multiFileOut bool,
	outputFile,
	outputFormat string) (ResultsWriter, error) {

	var resultsFormatter formatters.ResultsFormatter

	switch outputFormat {
	case "table":
		resultsFormatter = formatters.NewTableExtractionResultsFormatter()
	case "csv":
		resultsFormatter = formatters.NewCSVExtractionResultsFormatter()
	case "json":
		resultsFormatter = formatters.NewJsonExtractionResultsFormatter()
	}

	fileExtension := "." + outputFormat

	return newResultsWriter(multiFileOut, outputFile, fileExtension, resultsFormatter)
}

// NewReaderResultsWriter constructs a ResultsWriter configured for extraction.
func NewReaderResultsWriter(multiFileOut bool,
	outputFile, outputFormat string) (ResultsWriter, error) {
	fileExtension := ".ocr." + outputFormat

	var resultsFormatter formatters.ResultsFormatter = formatters.NewReadResultsFormatter()

	return newResultsWriter(multiFileOut, outputFile, fileExtension, resultsFormatter)
}

func newResultsWriter(multiFileOut bool, outputFile, fileExtension string,
	resultsFormatter formatters.ResultsFormatter) (ResultsWriter, error) {
	var resultsWriter ResultsWriter

	if multiFileOut {
		sinkFactory := sinks.NewExtensionSwappingFileSinkFactory(fileExtension)

		resultsWriter = NewIndividualResultsWriter(sinkFactory, resultsFormatter)
	} else {
		outFile, err := fs.OpenForWriting(outputFile)

		if err != nil {
			return nil, err
		}
		sink := sinks.NewBasicWriterSink(outFile)
		resultsWriter = NewCombinedResultsWriter(sink, resultsFormatter)
	}

	return resultsWriter, nil
}
