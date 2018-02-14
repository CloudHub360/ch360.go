package formatters

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/pkg/errors"
	"io"
	"path/filepath"
)

var _ ResultsFormatter = (*TableClassifyResultsFormatter)(nil)

type TableClassifyResultsFormatter struct {
}

func NewTableClassifyResultsFormatter() *TableClassifyResultsFormatter {
	return &TableClassifyResultsFormatter{}
}

var TableClassifyFormatterOutputFormat = "%-36.36s %-32.32s %v\n"

func (f *TableClassifyResultsFormatter) WriteResult(writer io.Writer, fullPath string, result interface{}, options FormatOption) error {

	classificationResult, ok := result.(*results.ClassificationResult)

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected type: %T", result))
	}

	if options&IncludeHeader == IncludeHeader {
		fmt.Fprintf(writer, TableClassifyFormatterOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	}

	filename := filepath.Base(fullPath)
	fmt.Fprintf(writer,
		TableClassifyFormatterOutputFormat,
		filename,
		classificationResult.DocumentType,
		classificationResult.IsConfident)

	return nil
}

func (f *TableClassifyResultsFormatter) Flush(writer io.Writer) error {
	return nil
}

func (f *TableClassifyResultsFormatter) Format() OutputFormat {
	return Table
}
