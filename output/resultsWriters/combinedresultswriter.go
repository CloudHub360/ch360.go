package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
)

// The CombinedResultsWriter writes all the results it is passed to a single resultSink (for either the console, or a single file),
// with appropriate header, record separators & footer provided by the specified resultsFormatter
type CombinedResultsWriter struct {
	writerFactory    sinks.SinkFactory
	resultsFormatter formatters.ClassifyResultsFormatter
	resultSink       sinks.Sink
	resultWritten    bool
}

func NewCombinedResultsWriter(sink sinks.Sink, resultsFormatter formatters.ClassifyResultsFormatter) *CombinedResultsWriter {
	return &CombinedResultsWriter{
		resultSink:       sink,
		resultsFormatter: resultsFormatter,
	}
}

func (c *CombinedResultsWriter) Start() error {
	c.resultSink.Open()
	c.resultsFormatter.WriteHeader(c.resultSink)
	return nil
}

func (c *CombinedResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	if c.resultWritten {
		c.resultsFormatter.WriteSeparator(c.resultSink)
	}

	c.resultsFormatter.WriteResult(c.resultSink, filename, result)
	c.resultWritten = true
	return nil
}

func (c *CombinedResultsWriter) Finish() error {
	c.resultsFormatter.WriteFooter(c.resultSink)

	c.resultSink.Close()
	return nil
}
