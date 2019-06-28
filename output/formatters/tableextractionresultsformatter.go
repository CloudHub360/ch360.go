package formatters

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"io"
	"path/filepath"
	"strings"
)

var _ ResultsFormatter = (*TableExtractionResultsFormatter)(nil)

type TableExtractionResultsFormatter struct {
}

func NewTableExtractionResultsFormatter() *TableExtractionResultsFormatter {
	return &TableExtractionResultsFormatter{}
}

var FieldColumnWidth uint = 31
var FileColumnWidth uint = 35

func (f *TableExtractionResultsFormatter) writeHeaderFor(writer io.Writer, result *results.ExtractionResult) error {
	cells := []*uitable.Cell{
		{Width: FileColumnWidth, Wrap: false, RightAlign: false, Data: "File"},
	}
	for _, fieldResult := range result.FieldResults {
		cells = append(cells, &uitable.Cell{
			Data:  fieldResult.FieldName,
			Width: FieldColumnWidth,
		})
	}

	row := uitable.Row{
		Cells:     cells,
		Separator: " ",
	}

	_, err := fmt.Fprintln(writer, strings.TrimSpace(row.String()))

	return err
}

func (f *TableExtractionResultsFormatter) WriteResult(writer io.Writer, fullPath string, result interface{}, options FormatOption) error {
	extractionResult, ok := result.(*results.ExtractionResult)

	if !ok {
		return errors.Errorf("unexpected type: %T", result)
	}

	if options&IncludeHeader == IncludeHeader {
		err := f.writeHeaderFor(writer, extractionResult)

		if err != nil {
			return err
		}
	}

	filename := filepath.Base(fullPath)

	row := uitable.Row{
		Separator: " ",
		Cells: []*uitable.Cell{
			{
				Wrap:  false,
				Width: FileColumnWidth,
				Data:  filename,
			},
		}}

	for _, fieldResult := range extractionResult.FieldResults {

		fieldFormatter := NewFieldFormatter(fieldResult, ", ", "(no result)")

		row.Cells = append(row.Cells, &uitable.Cell{
			Data:  fieldFormatter.String(),
			Width: FieldColumnWidth,
			Wrap:  true,
		})
	}

	_, err := fmt.Fprintln(writer, strings.TrimSpace(row.String()))

	return err
}

func (f *TableExtractionResultsFormatter) Flush(writer io.Writer) error {
	return nil
}
