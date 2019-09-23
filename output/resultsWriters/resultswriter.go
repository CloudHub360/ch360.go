package resultsWriters

import (
	"github.com/waives/surf/fs"
	"github.com/waives/surf/output/formatters"
	"github.com/waives/surf/output/sinks"
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

// NewClassificationResultsWriter constructs a ResultsWriter configured for classification.
func NewClassificationResultsWriter(multiFileOut bool,
	outputFile,
	outputFormat string) (ResultsWriter, error) {

	var resultsFormatter formatters.ResultsFormatter

	switch outputFormat {
	case "table":
		resultsFormatter = formatters.NewTableClassifyResultsFormatter()
	case "csv":
		resultsFormatter = formatters.NewCSVClassifyResultsFormatter()
	case "json":
		resultsFormatter = formatters.NewJsonClassifyResultsFormatter()
	}

	fileExtension := "." + outputFormat

	return newResultsWriter(multiFileOut, outputFile, fileExtension, resultsFormatter)
}

// NewReaderResultsWriter constructs a ResultsWriter configured for reading.
func NewReaderResultsWriter(multiFileOut bool,
	outputFile, outputFormat string) (ResultsWriter, error) {
	fileExtension := ".ocr." + outputFormat

	var resultsFormatter formatters.ResultsFormatter = formatters.NewNoopResultsFormatter()

	return newResultsWriter(multiFileOut, outputFile, fileExtension, resultsFormatter)
}

// NewReaderResultsWriter constructs a ResultsWriter configured for reading.
func NewRedactResultsWriter(multiFileOut bool,
	outputFile string) (ResultsWriter, error) {
	fileExtension := ".redacted.pdf"

	var resultsFormatter formatters.ResultsFormatter = formatters.NewNoopResultsFormatter()

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
