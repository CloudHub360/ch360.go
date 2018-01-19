package formatters

import (
	"encoding/csv"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"io"
	"path/filepath"
)

type CSVClassifyResultsFormatter struct {
}

func NewCSVClassifyResultsFormatter() *CSVClassifyResultsFormatter {
	return &CSVClassifyResultsFormatter{}
}

func (f *CSVClassifyResultsFormatter) WriteHeader(writer io.Writer) error {
	return nil
}

func (f *CSVClassifyResultsFormatter) WriteResult(writer io.Writer, filename string, result *types.ClassificationResult) error {
	record := []string{filepath.FromSlash(filename), result.DocumentType, boolToString(result.IsConfident), fmt.Sprintf("%.3f", result.RelativeConfidence)}

	csvWriter := csv.NewWriter(writer)

	if err := csvWriter.Write(record); err != nil {
		return err
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func (f *CSVClassifyResultsFormatter) WriteSeparator(writer io.Writer) error {
	return nil
}

func (f *CSVClassifyResultsFormatter) WriteFooter(writer io.Writer) error {
	return nil
}

func boolToString(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}
