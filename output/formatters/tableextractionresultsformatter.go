package formatters

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/pkg/errors"
	"io"
	"path/filepath"
)

var _ ResultsFormatter = (*TableClassifyResultsFormatter)(nil)

type TableExtractionResultsFormatter struct {
}

func NewTableExtractionResultsFormatter() *TableExtractionResultsFormatter {
	return &TableExtractionResultsFormatter{}
}

var FileHeaderFmt = "%-36.36s"
var FieldColFmt = "%-32.32s"

func (f *TableExtractionResultsFormatter) writeHeaderFor(writer io.Writer, result *results.ExtractionResult) error {
	fmt.Fprintf(writer, FileHeaderFmt, "FILE")

	for _, fieldResult := range result.FieldResults {
		fmt.Fprintf(writer, FieldColFmt, fieldResult.FieldName)
	}
	_, err := fmt.Fprintln(writer)
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

	fmt.Fprintf(writer,
		FileHeaderFmt,
		filename)

	for _, fieldResult := range extractionResult.FieldResults {
		fmt.Fprintf(writer, FieldColFmt, fieldResult.Result.Text)
	}

	fmt.Fprintln(writer)

	return nil
}

func (f *TableExtractionResultsFormatter) Flush(writer io.Writer) error {
	return nil
}

func (f *TableExtractionResultsFormatter) Format() OutputFormat {
	return Table
}
