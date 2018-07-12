package tests

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TableExtractionResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.TableExtractionResultsFormatter
	filename string
	result   *results.ExtractionResult
}

func (suite *TableExtractionResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewTableExtractionResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = anExtractionResult()
}

func TestTableExtractionResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(TableExtractionResultsFormatterSuite))
}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_FieldNames_In_Header() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, formatters.IncludeHeader)

	require.Nil(suite.T(), err)
	header := strings.Split(suite.output.String(), "\n")[0]
	suite.assertTableColumnsContent(header, []string{"File", "Amount"})
}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_FieldText() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	suite.assertTableColumnsContent(suite.output.String(), []string{suite.filename, "$5.50"})

}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_Empty_String_When_Field_Result_Is_Nil() {
	suite.result.FieldResults[0].Result = nil

	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	suite.assertTableColumnsContent(suite.output.String(), []string{suite.filename, formatters.NoResultText})
}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_Filename_Only_When_It_Has_Path() {
	// Arrange
	expectedFilename := `document1.tif`
	filename := filepath.Join(os.TempDir(), expectedFilename)

	// Act
	err := suite.sut.WriteResult(suite.output, filename, suite.result, 0)

	// Assert
	require.Nil(suite.T(), err)
	suite.assertTableColumnsContent(suite.output.String(), []string{expectedFilename, "$5.50"})
}

func (suite *TableExtractionResultsFormatterSuite) TestFlush_Writes_Nothing() {
	suite.sut.Flush(suite.output)

	assert.Equal(suite.T(), "", suite.output.String())
}

func (suite *TableExtractionResultsFormatterSuite) assertTableColumnsContent(row string, expectedContent []string) {
	var columns []string

	runes := []rune(row)

	var popResult = func(idx int) {
		columns = append(columns, strings.TrimSpace(string(runes[:idx])))
		runes = runes[idx:]
	}

	popResult(formatters.FileColumnWidth)

	for len(runes) > 0 {
		idx := formatters.FieldColumnWidth
		if len(runes) < idx {
			idx = len(runes)
		}
		popResult(idx)
	}

	suite.Assert().Equal(expectedContent, columns)
}
