package formatters

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"io"
	"path/filepath"
)

type CSVClassifyResultsFormatter struct {
}

var _ (ResultsFormatter) = (*CSVClassifyResultsFormatter)(nil)

func NewCSVClassifyResultsFormatter() *CSVClassifyResultsFormatter {
	return &CSVClassifyResultsFormatter{}
}

func (f *CSVClassifyResultsFormatter) WriteResult(writer io.Writer, filename string, result interface{}, options FormatOption) error {
	classificationResult, ok := result.(*results.ClassificationResult)

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected type: %T", result))
	}

	csvWriter := csv.NewWriter(writer)

	if options&IncludeHeader == IncludeHeader {
		csvWriter.Write([]string{"File", "Document Type", "Confident", "Relative Confidence"})
	}

	record := []string{filepath.FromSlash(filename),
		classificationResult.DocumentType,
		fmt.Sprintf("%v", classificationResult.IsConfident),
		fmt.Sprintf("%.3f", classificationResult.RelativeConfidence)}

	err := csvWriter.Write(record)

	if err != nil {
		return err
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return err
	}

	return nil
}

func (f *CSVClassifyResultsFormatter) Flush(writer io.Writer) error {
	return nil
}

func (f *CSVClassifyResultsFormatter) Format() OutputFormat {
	return Csv
}
