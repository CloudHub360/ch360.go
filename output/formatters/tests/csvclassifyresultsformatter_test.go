package tests

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/output/formatters"
	"github.com/waives/surf/test/generators"
	"strings"
	"testing"
)

type CSVClassifyResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.CSVClassifyResultsFormatter
	filename string
	result   *results.ClassificationResult
}

func (suite *CSVClassifyResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewCSVClassifyResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = &results.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestCSVResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(CSVClassifyResultsFormatterSuite))
}

func (suite *CSVClassifyResultsFormatterSuite) TestWrites_ResultWithCorrectFormat() {
	filename := "document1.tif"
	result := &results.ClassificationResult{
		DocumentType:       "documenttype",
		IsConfident:        true,
		RelativeConfidence: 1.234567,
	}
	expectedOutput := "document1.tif,documenttype,true,1.235\n"

	err := suite.sut.WriteResult(suite.output, filename, result, 0)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedOutput, suite.output.String())
}

func (suite *CSVClassifyResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *CSVClassifyResultsFormatterSuite) TestWrites_DocumentType() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *CSVClassifyResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *CSVClassifyResultsFormatterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
	filename := `C:\\folder\\document1.tif`

	err := suite.sut.WriteResult(suite.output, filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), filename, suite.output.String()[:len(filename)])
}

func (suite *CSVClassifyResultsFormatterSuite) TestFlush_Writes_Nothing() {
	suite.sut.Flush(suite.output)

	assert.Equal(suite.T(), "", suite.output.String())
}
