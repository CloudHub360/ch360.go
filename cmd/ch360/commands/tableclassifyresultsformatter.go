package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"path/filepath"
)

//go:generate mockery -name "ClassifyResultsFormatter"
type ClassifyResultsFormatter interface {
	Start() error
	WriteResult(filename string, result *types.ClassificationResult) error
	Finish() error
}

type TableClassifyResultsFormatter struct {
	provider WriterProvider
}

func NewTableClassifyResultsFormatter(provider WriterProvider) *TableClassifyResultsFormatter {
	return &TableClassifyResultsFormatter{
		provider: provider,
	}
}

var classifyTableOutputFormat = "%-44.44s %-24.24s %v\n"

func (writer *TableClassifyResultsFormatter) Start() error {
	outWriter, err := writer.provider.Provide("")

	if err != nil {
		return err
	}

	if outWriter != nil {
		fmt.Fprintf(outWriter, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	}

	return nil
}

func (writer *TableClassifyResultsFormatter) WriteResult(fullPath string, result *types.ClassificationResult) error {
	out, err := writer.provider.Provide(fullPath)

	if err != nil {
		return err
	}

	filename := filepath.Base(fullPath)
	fmt.Fprintf(out, ClassifyOutputFormat, filename, result.DocumentType, result.IsConfident)

	return nil
}

func (writer *TableClassifyResultsFormatter) Finish() error {
	return nil
}
