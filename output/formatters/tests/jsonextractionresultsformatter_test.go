package tests

import (
	"bytes"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"runtime"
	"strings"
	"testing"
)

type JsonExtractionResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.JsonExtractionResultsFormatter
	filename string
	result   *results.ExtractionResult
}

func (suite *JsonExtractionResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewJsonExtractionResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = anExtractionResult()
}

func TestJsonExtractionResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(JsonExtractionResultsFormatterSuite))
}

func (suite *JsonExtractionResultsFormatterSuite) TestWriteResult_With_IncludeHeader_Option_Writes_Header() {
	suite.sut.WriteResult(suite.output, suite.filename, suite.result, formatters.IncludeHeader)

	assert.True(suite.T(), strings.HasPrefix(suite.output.String(), "["))
}

func (suite *JsonExtractionResultsFormatterSuite) TestWrites_ResultWithCorrectFormat_Without_Header() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)
	suite.sut.Flush(suite.output)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), IsJSON(suite.output.String()))
}

func (suite *JsonExtractionResultsFormatterSuite) TestWrites_ResultWithCorrectFormat_With_Header() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, formatters.IncludeHeader)
	suite.sut.Flush(suite.output)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), IsJSON(suite.output.String()))
	assert.True(suite.T(), strings.HasPrefix(suite.output.String(), "["))
	assert.True(suite.T(), strings.HasSuffix(suite.output.String(), "]"))
}

func (suite *JsonExtractionResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *JsonExtractionResultsFormatterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
	var filename string
	var expectedFilename string

	if runtime.GOOS == "windows" {
		filename = `C:/folder/document1.tif`
		expectedFilename = `C:\\folder\\document1.tif` //Json pretty-printing escapes slashes
	} else {
		filename = `/var/folder/document1.tif`
		expectedFilename = `/var/folder/document1.tif`
	}

	err := suite.sut.WriteResult(suite.output, filename, suite.result, 0)

	require.Nil(suite.T(), err)
	containsFilenameWithPath := strings.Contains(suite.output.String(), expectedFilename)

	assert.True(suite.T(), containsFilenameWithPath)
	if !containsFilenameWithPath { // To aid debugging if test fails
		fmt.Println(suite.output.String())
		fmt.Printf("Output does not contain %s", expectedFilename)
	}
}

func (suite *JsonExtractionResultsFormatterSuite) TestResult_Param_Must_Be_Correct_Type() {
	err := suite.sut.WriteResult(suite.output, suite.filename, struct{}{}, 0)

	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.HasPrefix(err.Error(), "Unexpected type"))
}

func (suite *JsonExtractionResultsFormatterSuite) TestSeparator_Written_When_Outputting_Multiple_Results() {
	suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)
	suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)
	suite.sut.Flush(suite.output)

	assert.True(suite.T(), IsJSON(suite.output.String()))
}
