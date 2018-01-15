package commands

import (
	"fmt"
	"io"
	"path/filepath"
)

type ClassifyResultsWriter interface {
	StartWriting()
	WriteDocumentResults(result *classifyResultsWriterInput)
	FinishWriting()
}

type classifyResultsWriterInput struct {
	filename     string //Filename with full path
	documentType string
	isConfident  bool
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

func (writer *TableClassifyResultsWriter) WriteDocumentResults(result *classifyResultsWriterInput) {
	base := filepath.Base(result.filename)
	fmt.Fprintf(writer.writer, ClassifyOutputFormat, base, result.documentType, result.isConfident)
}

func (writer *TableClassifyResultsWriter) FinishWriting() {}
