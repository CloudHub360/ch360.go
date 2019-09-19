package tests

import (
	"bytes"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/output/formatters"
	"github.com/waives/surf/test/generators"
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
	suite.assertRecordFieldContent(suite.output.String(), []string{suite.filename, "$5.50"})
}

func (suite *CSVExtractionResultsFormatterSuite) TestWrites_Column_Per_Field() {
	expectedHeadings := []string{"Filename", "Amount", "Amount2"}

	err := suite.sut.WriteResult(suite.output, suite.filename, anExtractionResultWith2Fields(), formatters.IncludeHeader)

	require.Nil(suite.T(), err)
	require.True(suite.T(), isCSV(suite.output.String()))
	record, err := csvReaderFor(suite.output.String()).Read()
	assert.Equal(suite.T(), expectedHeadings, record)
}

func (suite *CSVExtractionResultsFormatterSuite) TestWrites_Header_Row_When_Specified() {
	expectedHeadings := []string{"Filename", "Amount"}
	expectedData := []string{suite.filename, "$5.50"}
	outputWithHeader := &bytes.Buffer{}
	outputWithoutHeader := &bytes.Buffer{}

	suite.sut.WriteResult(outputWithHeader, suite.filename, anExtractionResult(), formatters.IncludeHeader)
	suite.sut.WriteResult(outputWithoutHeader, suite.filename, anExtractionResult(), 0)

	suite.assertRecordFieldContent(outputWithHeader.String(), expectedHeadings)
	suite.assertRecordFieldContent(outputWithoutHeader.String(), expectedData)
}

func (suite *CSVExtractionResultsFormatterSuite) TestWrites_Empty_String_On_Null_Result() {
	result := anExtractionResult()
	result.FieldResults[0].Result = nil // remove field result
	expectedValues := []string{suite.filename, ""}

	err := suite.sut.WriteResult(suite.output, suite.filename, result, 0)

	require.Nil(suite.T(), err)
	suite.assertRecordFieldContent(suite.output.String(), expectedValues)
}

func (suite *CSVExtractionResultsFormatterSuite) assertRecordFieldContent(recordStr string, content []string) {
	csvReader := csvReaderFor(recordStr)
	record, _ := csvReader.Read()
	suite.Assert().Equal(content, record)
}

func isCSV(data string) bool {
	csvReader := csvReaderFor(data)
	_, err := csvReader.Read()

	return err == nil
}

func csvReaderFor(data string) *csv.Reader {
	return csv.NewReader(strings.NewReader(data))
}
