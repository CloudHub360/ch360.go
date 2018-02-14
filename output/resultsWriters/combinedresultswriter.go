package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
)

var _ ResultsWriter = (*CombinedResultsWriter)(nil)

// The CombinedResultsWriter writes all the results it is passed to a single resultSink (for either the console, or a single file),
// with appropriate header, record separators & footer provided by the specified resultsFormatter
type CombinedResultsWriter struct {
	resultsFormatter formatters.ResultsFormatter
	resultSink       sinks.Sink
	resultWritten    bool
}

func NewCombinedResultsWriter(sink sinks.Sink, resultsFormatter formatters.ResultsFormatter) *CombinedResultsWriter {
	return &CombinedResultsWriter{
		resultSink:       sink,
		resultsFormatter: resultsFormatter,
	}
}

func (c *CombinedResultsWriter) Start() error {
	if err := c.resultSink.Open(); err != nil {
		return err
	}

	return nil
}

func (c *CombinedResultsWriter) WriteResult(filename string, result interface{}) error {
	var formatOptions formatters.FormatOption = 0
	if !c.resultWritten {
		formatOptions = formatters.IncludeHeader
	}

	c.resultWritten = true
	return c.resultsFormatter.WriteResult(c.resultSink, filename, result, formatOptions)
}

func (c *CombinedResultsWriter) Finish() error {
	err := c.resultsFormatter.Flush(c.resultSink)

	if err != nil {
		return err
	}

	return c.resultSink.Close()
}
