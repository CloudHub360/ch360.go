package commands

import (
	"encoding/csv"
	"io"
)

type CSVClassifyResultsWriter struct {
	underlyingWriter io.Writer
	csvWriter        *csv.Writer
}

func NewCSVClassifyResultsWriter(writer io.Writer) *CSVClassifyResultsWriter {
	return &CSVClassifyResultsWriter{
		underlyingWriter: writer,
	}
}

func (writer *CSVClassifyResultsWriter) StartWriting() {
	writer.csvWriter = csv.NewWriter(writer.underlyingWriter)
}

func (writer *CSVClassifyResultsWriter) WriteDocumentResults(result *classifyResultsWriterInput) {
	record := []string{result.filename, result.documentType, boolToString(result.isConfident)}

	writer.csvWriter.Write(record) //TODO check error
	writer.csvWriter.Flush()
}

func (writer *CSVClassifyResultsWriter) FinishWriting() {}

func boolToString(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}
