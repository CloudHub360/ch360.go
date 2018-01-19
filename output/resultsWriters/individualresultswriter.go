package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
)

// The IndividualResultsWriter writes each results it is passed to a new resultSink (an extensionSwappingFileWriter),
// just using the WriteResult method provided by the formatter (no header, record separators or footer)
type IndividualResultsWriter struct {
	sinkFactory      sinks.SinkFactory
	resultsFormatter formatters.ClassifyResultsFormatter
}

func NewIndividualResultsWriter(writerFactory sinks.SinkFactory, resultsFormatter formatters.ClassifyResultsFormatter) *IndividualResultsWriter {
	return &IndividualResultsWriter{
		sinkFactory:      writerFactory,
		resultsFormatter: resultsFormatter,
	}
}

func (c *IndividualResultsWriter) Start() error {
	return nil
}

func (c *IndividualResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	resultSink, _ := c.sinkFactory.Sink(
		sinks.SinkParams{InputFilename: filename}) //Returns resultSink configured to write to destination file

	resultSink.Open()
	c.resultsFormatter.WriteResult(resultSink, filename, result)

	resultSink.Close()
	return nil
}

func (c *IndividualResultsWriter) Finish() error {
	return nil
}
