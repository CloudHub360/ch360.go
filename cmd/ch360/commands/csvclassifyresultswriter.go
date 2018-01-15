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

func (writer *CSVClassifyResultsWriter) WriteDocumentResults(result *classifyResultsWriterInput) error {
	record := []string{result.filename, result.documentType, boolToString(result.isConfident)}

	if err := writer.csvWriter.Write(record); err != nil {
		return err
	}

	writer.csvWriter.Flush()

	if err := writer.csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func (writer *CSVClassifyResultsWriter) FinishWriting() {}

func boolToString(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}
