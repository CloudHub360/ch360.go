package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
)

var _ ResultsWriter = (*IndividualResultsWriter)(nil)

// The IndividualResultsWriter writes each results it is passed to a new resultSink (an extensionSwappingFileWriter),
// just using the WriteResult method provided by the formatter (no header, record separators or footer)
type IndividualResultsWriter struct {
	sinkFactory      sinks.SinkFactory
	resultsFormatter formatters.ResultsFormatter
}

func NewIndividualResultsWriter(writerFactory sinks.SinkFactory, resultsFormatter formatters.ResultsFormatter) *IndividualResultsWriter {
	return &IndividualResultsWriter{
		sinkFactory:      writerFactory,
		resultsFormatter: resultsFormatter,
	}
}

func (c *IndividualResultsWriter) Start() error {
	return nil
}

func (c *IndividualResultsWriter) WriteResult(filename string, result interface{}) error {
	resultSink, err := c.sinkFactory.Sink(
		sinks.SinkParams{InputFilename: filename}) //Returns resultSink configured to write to destination file

	if err != nil {
		return err
	}

	err = resultSink.Open()
	if err != nil {
		return err
	}

	err = c.resultsFormatter.WriteResult(resultSink, filename, result, formatters.IncludeHeader)
	if err != nil {
		return err
	}

	err = c.resultsFormatter.Flush(resultSink)
	if err != nil {
		return err
	}

	err = resultSink.Close()

	if err != nil {
		return err
	}

	return nil
}

func (c *IndividualResultsWriter) Finish() error {
	return nil
}
