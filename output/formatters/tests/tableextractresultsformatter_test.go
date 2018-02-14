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

func (suite *TableExtractionResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_FieldNames_In_Header() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, formatters.IncludeHeader)

	require.Nil(suite.T(), err)
	for _, field := range suite.result.FieldResults {
		assert.True(suite.T(), strings.Contains(suite.output.String(), field.FieldName))
	}
}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_FieldText() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	for _, field := range suite.result.FieldResults {
		assert.True(suite.T(), strings.Contains(suite.output.String(), field.Result.Text))
	}
}

func (suite *TableExtractionResultsFormatterSuite) TestWrites_Filename_Only_When_It_Has_Path() {
	// Arrange
	expectedFilename := `document1.tif`
	filename := filepath.Join(os.TempDir(), expectedFilename)

	// Act
	err := suite.sut.WriteResult(suite.output, filename, suite.result, 0)

	// Assert
	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), expectedFilename))
}

func (suite *TableExtractionResultsFormatterSuite) TestFlush_Writes_Nothing() {
	suite.sut.Flush(suite.output)

	assert.Equal(suite.T(), "", suite.output.String())
}
