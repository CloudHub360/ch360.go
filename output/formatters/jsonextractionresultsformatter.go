package formatters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"io"
)

type JsonExtractionResultsFormatter struct {
	resultsWritten bool
	headerWritten  bool
}

var _ ResultsFormatter = (*JsonExtractionResultsFormatter)(nil)

func NewJsonExtractionResultsFormatter() *JsonExtractionResultsFormatter {
	return &JsonExtractionResultsFormatter{}
}

func (f *JsonExtractionResultsFormatter) WriteResult(writer io.Writer, filename string, result interface{}, options FormatOption) error {

	extractionResult, ok := result.(*results.ExtractionResult)

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected type: %T", result))
	}

	if options&IncludeHeader == IncludeHeader || !f.resultsWritten {
		f.headerWritten = true
		fmt.Fprint(writer, "[") // header
	} else if f.resultsWritten {
		fmt.Fprint(writer, ",\n") // write separator
	}

	// add filename to original result
	var output = struct {
		Filename string `json:"filename"`
		*results.ExtractionResult
	}{
		Filename:         filename,
		ExtractionResult: extractionResult,
	}

	bytes, err := json.MarshalIndent(&output, "", "  ")

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)

	f.resultsWritten = true

	return err
}

func (f *JsonExtractionResultsFormatter) Flush(writer io.Writer) error {
	if f.headerWritten {
		_, err := fmt.Fprint(writer, "]")
		return err
	}
	return nil
}

func (f *JsonExtractionResultsFormatter) Format() OutputFormat {
	return Json
}
