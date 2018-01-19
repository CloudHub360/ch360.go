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
	if err := c.resultSink.Open(); err != nil {
		return err
	}

	return c.resultsFormatter.WriteHeader(c.resultSink)
}

func (c *CombinedResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	if c.resultWritten {
		if err := c.resultsFormatter.WriteSeparator(c.resultSink); err != nil {
			return err
		}
	}

	err := c.resultsFormatter.WriteResult(c.resultSink, filename, result)
	c.resultWritten = true
	return err
}

func (c *CombinedResultsWriter) Finish() error {
	if err := c.resultsFormatter.WriteFooter(c.resultSink); err != nil {
		return err
	}

	return c.resultSink.Close()
}
