package formatters

import (
	"encoding/csv"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/pkg/errors"
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
		return errors.Errorf("unexpected type: %T", result)
	}

	if options&IncludeHeader == IncludeHeader {
		f.writeHeaderFor(writer, extractionResult)
	}

	record := []string{filepath.FromSlash(filename)}

	for _, fieldResult := range extractionResult.FieldResults {

		fieldFormatter := NewFieldFormatter(fieldResult, "|", "")

		record = append(record, fieldFormatter.String())
	}

	return f.writeRecord(writer, record)
}

func (f *CSVExtractionResultsFormatter) Flush(writer io.Writer) error {
	return nil
}
