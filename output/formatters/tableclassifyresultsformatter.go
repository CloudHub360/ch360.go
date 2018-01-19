package formatters

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"io"
	"path/filepath"
)

//go:generate mockery -name "ClassifyResultsFormatter"
type ClassifyResultsFormatter interface {
	WriteHeader(writer io.Writer) error
	WriteResult(writer io.Writer, filename string, result *types.ClassificationResult) error
	WriteSeparator(writer io.Writer) error
	WriteFooter(writer io.Writer) error
}

type TableClassifyResultsFormatter struct {
}

func NewTableClassifyResultsFormatter() *TableClassifyResultsFormatter {
	return &TableClassifyResultsFormatter{}
}

var TableFormatterOutputFormat = "%-36.36s %-32.32s %v\n"

func (f *TableClassifyResultsFormatter) WriteHeader(writer io.Writer) error {
	if writer != nil {
		fmt.Fprintf(writer, TableFormatterOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	}

	return nil
}

func (f *TableClassifyResultsFormatter) WriteResult(writer io.Writer, fullPath string, result *types.ClassificationResult) error {
	filename := filepath.Base(fullPath)
	fmt.Fprintf(writer, TableFormatterOutputFormat, filename, result.DocumentType, result.IsConfident)

	return nil
}

func (f *TableClassifyResultsFormatter) WriteSeparator(writer io.Writer) error {
	return nil
}

func (f *TableClassifyResultsFormatter) WriteFooter(writer io.Writer) error {
	return nil
}
