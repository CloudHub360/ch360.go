package formatters

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"io"
	"path/filepath"
)

var _ (ResultsFormatter) = (*CSVExtractionResultsFormatter)(nil)

type CSVExtractionResultsFormatter struct {
}

func NewCSVExtractionResultsFormatter() *CSVExtractionResultsFormatter {
	return &CSVExtractionResultsFormatter{}
}

func (f *CSVExtractionResultsFormatter) writeHeaderFor(writer io.Writer, result *results.ExtractionResult) error {

	record := []string{"Filename"}

	for _, field := range result.FieldResults {
		record = append(record, field.FieldName)
	}

	return f.writeRecord(writer, record)

}

func (f *CSVExtractionResultsFormatter) writeRecord(writer io.Writer, record []string) error {
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

func (f *CSVExtractionResultsFormatter) WriteResult(writer io.Writer, filename string, result interface{}, options FormatOption) error {
	extractionResult, ok := result.(*results.ExtractionResult)

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected type: %T", result))
	}

	if options&IncludeHeader == IncludeHeader {
		f.writeHeaderFor(writer, extractionResult)
	}

	record := []string{filepath.FromSlash(filename)}

	for _, field := range extractionResult.FieldResults {
		resultText := ""
		if field.Result != nil {
			resultText = field.Result.Text
		}
		record = append(record, resultText)
	}

	return f.writeRecord(writer, record)
}

func (f *CSVExtractionResultsFormatter) Flush(writer io.Writer) error {
	return nil
}
