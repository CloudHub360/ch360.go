package commands

import (
	"encoding/csv"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"path/filepath"
)

type CSVClassifyResultsWriter struct {
	writerProvider WriterProvider
}

func NewCSVClassifyResultsWriter(provider WriterProvider) *CSVClassifyResultsWriter {
	return &CSVClassifyResultsWriter{
		writerProvider: provider,
	}
}

func (writer *CSVClassifyResultsWriter) Start() error {
	return nil
}

func (writer *CSVClassifyResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	destWriter, err := writer.writerProvider.Provide(filename)

	if err != nil {
		return err
	}

	record := []string{filepath.FromSlash(filename), result.DocumentType, boolToString(result.IsConfident), fmt.Sprintf("%.3f", result.RelativeConfidence)}

	csvWriter := csv.NewWriter(destWriter)

	if err := csvWriter.Write(record); err != nil {
		return err
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func (writer *CSVClassifyResultsWriter) Finish() error {
	return nil
}

func boolToString(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}
