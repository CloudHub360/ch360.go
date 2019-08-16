package formatters

import (
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
