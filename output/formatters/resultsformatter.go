package formatters

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"io"
)

type OutputFormat string

const (
	Table OutputFormat = "table"
	Json  OutputFormat = "json"
	Csv   OutputFormat = "csv"
)

type FormatOption int

const (
	IncludeHeader FormatOption = 1 << iota
	IncludeSeparator
)

//go:generate mockery -name "ResultsFormatter"
type ResultsFormatter interface {
	WriteResult(writer io.Writer, filename string, result interface{}, options FormatOption) error
	Flush(writer io.Writer) error
}

func NewResultsFormatterFor(params *config.GlobalFlags, verb config.Verb) (ResultsFormatter, error) {
	if verb == config.Read {
		return NewReadResultsFormatter(), nil
	}

	var formatter ResultsFormatter
	switch OutputFormat(params.OutputFormat) {
	case Table:
		if verb == config.Classify {
			formatter = NewTableClassifyResultsFormatter()
		} else {
			formatter = NewTableExtractionResultsFormatter()
		}
	case Csv:
		if verb == config.Classify {
			formatter = NewCSVClassifyResultsFormatter()
		} else {
			formatter = NewCSVExtractionResultsFormatter()
		}
	case Json:
		if verb == config.Classify {
			formatter = NewJsonClassifyResultsFormatter()
		} else {
			formatter = NewJsonExtractionResultsFormatter()
		}
	default:
		// DocOpt doesn't do validation of these values for us, so we need to catch invalid values here
		return nil, errors.New(fmt.Sprintf("Unknown output format '%s'. Allowed values are: csv, table, json.", params.OutputFormat))
	}

	return formatter, nil
}
