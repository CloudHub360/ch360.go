package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"path/filepath"
)

//go:generate mockery -name "ClassifyResultsWriter"
type ClassifyResultsWriter interface {
	Start() error
	WriteResult(filename string, result *types.ClassificationResult) error
	Finish() error
}

type TableClassifyResultsWriter struct {
	provider WriteCloserProvider
}

func NewTableClassifyResultsWriter(provider WriteCloserProvider) *TableClassifyResultsWriter {
	return &TableClassifyResultsWriter{
		provider: provider,
	}
}

var classifyTableOutputFormat = "%-44.44s %-24.24s %v\n"

func (writer *TableClassifyResultsWriter) Start() error {
	outWriter, err := writer.provider("")

	if err != nil {
		return err
	}

	if outWriter != nil {
		fmt.Fprintf(outWriter, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	}

	return nil
}

func (writer *TableClassifyResultsWriter) WriteResult(fullPath string, result *types.ClassificationResult) error {
	out, err := writer.provider(fullPath)

	if err != nil {
		return err
	}

	filename := filepath.Base(fullPath)
	fmt.Fprintf(out, ClassifyOutputFormat, filename, result.DocumentType, result.IsConfident)

	return nil
}

func (writer *TableClassifyResultsWriter) Finish() error {
	outWriter, err := writer.provider("")

	if err != nil {
		return err
	}

	if outWriter != nil {
		return outWriter.Close()
	}

	return nil
}
