package formatters

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/pkg/errors"
	"io"
	"path/filepath"
	"strings"
)

var _ ResultsFormatter = (*TableExtractionResultsFormatter)(nil)

type TableExtractionResultsFormatter struct {
}

func NewTableExtractionResultsFormatter() *TableExtractionResultsFormatter {
	return &TableExtractionResultsFormatter{}
}

var FileColumnWidth = 36
var FieldColumnWidth = 32
var FileColumnFmt = fmt.Sprintf("%%-%d.%ds", FileColumnWidth, FileColumnWidth)
var FieldColumnFmt = fmt.Sprintf("%%-%d.%ds", FieldColumnWidth, FieldColumnWidth)
var NoResultText = "(no result)"

func (f *TableExtractionResultsFormatter) writeHeaderFor(writer io.Writer, result *results.ExtractionResult) error {
	var header = fmt.Sprintf(FileColumnFmt, "File")

	for _, fieldResult := range result.FieldResults {
		header = header + fmt.Sprintf(FieldColumnFmt, fieldResult.FieldName)
	}
	_, err := fmt.Fprintln(writer, strings.TrimSpace(header))
	return err
}

func (f *TableExtractionResultsFormatter) WriteResult(writer io.Writer, fullPath string, result interface{}, options FormatOption) error {
	extractionResult, ok := result.(*results.ExtractionResult)

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected type: %T", result))
	}

	if options&IncludeHeader == IncludeHeader {
		f.writeHeaderFor(writer, extractionResult)
	}

	filename := filepath.Base(fullPath)

	var row = fmt.Sprintf(FileColumnFmt, filename)

	for _, fieldResult := range extractionResult.FieldResults {

		resultText := FieldFormatter{FieldResult: fieldResult}.String()

		row = row + fmt.Sprintf(FieldColumnFmt, resultText)
	}

	fmt.Fprintln(writer, strings.TrimSpace(row))

	return nil
}

func (f *TableExtractionResultsFormatter) Flush(writer io.Writer) error {
	return nil
}
