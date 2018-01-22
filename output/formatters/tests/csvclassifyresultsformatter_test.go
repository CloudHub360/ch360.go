package tests

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type CSVResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.CSVClassifyResultsFormatter
	filename string
	result   *types.ClassificationResult
}

func (suite *CSVResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewCSVClassifyResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = &types.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestCSVResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(CSVResultsFormatterSuite))
}

func (suite *CSVResultsFormatterSuite) TestWriteHeader_Does_Not_Write_Anything() {
	suite.sut.WriteHeader(suite.output)

	assert.Equal(suite.T(), "", suite.output.String())
}

func (suite *CSVResultsFormatterSuite) TestWrites_ResultWithCorrectFormat() {
	filename := "document1.tif"
	result := &types.ClassificationResult{
		DocumentType:       "documenttype",
		IsConfident:        true,
		RelativeConfidence: 1.234567,
	}
	expectedOutput := "document1.tif,documenttype,true,1.235\n"

	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, filename, result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedOutput, suite.output.String())
}

func (suite *CSVResultsFormatterSuite) TestWrites_Filename() {
	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *CSVResultsFormatterSuite) TestWrites_DocumentType() {
	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *CSVResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *CSVResultsFormatterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
	filename := `C:\folder\document1.tif`

	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, filename, suite.result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), filename, suite.output.String()[:len(filename)])
}

//TODO Test that WriteFooter writes nothing
