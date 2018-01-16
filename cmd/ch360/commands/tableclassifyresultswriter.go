package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"io"
	"path/filepath"
)

type ClassifyResultsWriter interface {
	StartWriting()
	WriteDocumentResults(filename string, result *types.ClassificationResult) error
	FinishWriting()
}

type TableClassifyResultsWriter struct {
	writer io.Writer
}

func NewTableClassifyResultsWriter(writer io.Writer) *TableClassifyResultsWriter {
	return &TableClassifyResultsWriter{
		writer: writer,
	}
}

var classifyTableOutputFormat = "%-44.44s %-24.24s %v\n"

func (writer *TableClassifyResultsWriter) StartWriting() {
	fmt.Fprintf(writer.writer, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
}

func (writer *TableClassifyResultsWriter) WriteDocumentResults(filename string, result *types.ClassificationResult) error {
	base := filepath.Base(filename)
	fmt.Fprintf(writer.writer, ClassifyOutputFormat, base, result.DocumentType, result.IsConfident)

	return nil
}

func (writer *TableClassifyResultsWriter) FinishWriting() {}
