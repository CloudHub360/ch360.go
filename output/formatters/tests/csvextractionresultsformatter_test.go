package tests

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type CSVExtractionResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.CSVExtractionResultsFormatter
	filename string
	result   *results.ExtractionResult
}

func (suite *CSVExtractionResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewCSVExtractionResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = anExtractionResult()
}

func TestCSVExtractionResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(CSVExtractionResultsFormatterSuite))
}

func (suite *CSVExtractionResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

//
//func (suite *CSVExtractionResultsFormatterSuite) TestWrites_DocumentType() {
//	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)
//
//	require.Nil(suite.T(), err)
//	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
//}
//
//func (suite *CSVExtractionResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
//	suite.result.IsConfident = false
//
//	suite.sut.WriteHeader(suite.output)
//	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)
//
//	require.Nil(suite.T(), err)
//	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
//}
//
//func (suite *CSVExtractionResultsFormatterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
//	filename := `C:\folder\document1.tif`
//
//	suite.sut.WriteHeader(suite.output)
//	err := suite.sut.WriteResult(suite.output, filename, suite.result)
//
//	require.Nil(suite.T(), err)
//	assert.Equal(suite.T(), filename, suite.output.String()[:len(filename)])
//}
//
//func (suite *CSVExtractionResultsFormatterSuite) TestWriteFooter_Writes_Nothing() {
//	suite.sut.WriteFooter(suite.output)
//
//	assert.Equal(suite.T(), "", suite.output.String())
//}
//
//func (suite *CSVExtractionResultsFormatterSuite) TestWriteHeader_Writes_Nothing() {
//	suite.sut.WriteHeader(suite.output)
//
//	assert.Equal(suite.T(), "", suite.output.String())
//}
