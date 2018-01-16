package commands

import (
	"encoding/csv"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/pkg/errors"
	"io"
)

type CSVClassifyResultsWriter struct {
	underlyingWriter io.Writer
	csvWriter        *csv.Writer
	startCalled      bool
}

func NewCSVClassifyResultsWriter(writer io.Writer) *CSVClassifyResultsWriter {
	return &CSVClassifyResultsWriter{
		underlyingWriter: writer,
	}
}

func (writer *CSVClassifyResultsWriter) Start() {
	writer.csvWriter = csv.NewWriter(writer.underlyingWriter)
	writer.startCalled = true
}

func (writer *CSVClassifyResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	if !writer.startCalled {
		return errors.New("Start() must be called before WriteResult()")
	}

	record := []string{filename, result.DocumentType, boolToString(result.IsConfident), fmt.Sprintf("%.3f", result.RelativeConfidence)}

	if err := writer.csvWriter.Write(record); err != nil {
		return err
	}

	writer.csvWriter.Flush()

	if err := writer.csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func (writer *CSVClassifyResultsWriter) Finish() {}

func boolToString(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}
