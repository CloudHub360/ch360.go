package commands

import (
	"fmt"
	"io"
)

type CSVClassifyResultsWriter struct {
	writer io.Writer
}

func NewCSVClassifyResultsWriter(writer io.Writer) *CSVClassifyResultsWriter {
	return &CSVClassifyResultsWriter{
		writer: writer,
	}
}

var classifyCSVOutputFormat = "%s,%s,%v\n"

func (writer *CSVClassifyResultsWriter) StartWriting() {}

func (writer *CSVClassifyResultsWriter) WriteDocumentResults(result *classifyResultsWriterInput) {
	fmt.Fprintf(writer.writer, classifyCSVOutputFormat, result.filename, result.documentType, result.isConfident)
}

func (writer *CSVClassifyResultsWriter) FinishWriting() {}
