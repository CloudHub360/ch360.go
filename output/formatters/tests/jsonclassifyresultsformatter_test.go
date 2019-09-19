package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/output/formatters"
	"github.com/waives/surf/test/generators"
	"runtime"
	"strings"
	"testing"
)

type JsonClassificationResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.JsonClassifyResultsFormatter
	filename string
	result   *results.ClassificationResult
}

func (suite *JsonClassificationResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewJsonClassifyResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = &results.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestJsonClassificationResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(JsonClassificationResultsFormatterSuite))
}

func (suite *JsonClassificationResultsFormatterSuite) TestWriteResult_With_IncludeHeader_Option_Writes_Header() {
	suite.sut.WriteResult(suite.output, exampleFilename, exampleResult, formatters.IncludeHeader)

	assert.True(suite.T(), strings.HasPrefix(suite.output.String(), "["))
}

func (suite *JsonClassificationResultsFormatterSuite) TestWrites_ResultWithCorrectFormat_Without_Header() {
	err := suite.sut.WriteResult(suite.output, exampleFilename, exampleResult, 0)
	suite.sut.Flush(suite.output)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), exampleOutputObject, suite.output.String())
}

func (suite *JsonClassificationResultsFormatterSuite) TestWrites_ResultWithCorrectFormat_With_Header() {
	err := suite.sut.WriteResult(suite.output, exampleFilename, exampleResult, formatters.IncludeHeader)
	suite.sut.Flush(suite.output)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), exampleOutputList, suite.output.String())
}

func (suite *JsonClassificationResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *JsonClassificationResultsFormatterSuite) TestWrites_DocumentType() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *JsonClassificationResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *JsonClassificationResultsFormatterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
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

var exampleFilename = "document1.tif"
var exampleResult = &results.ClassificationResult{
	DocumentType:       "documenttype",
	IsConfident:        true,
	RelativeConfidence: 1.234567,
	DocumentTypeScores: []results.DocumentTypeScore{
		results.DocumentTypeScore{
			DocumentType: "documenttype",
			Score:        1.23456,
		},
		results.DocumentTypeScore{
			DocumentType: "otherdocumenttype",
			Score:        33.45678,
		},
	}}

var exampleOutputObject = `{
  "filename": "document1.tif",
  "classification_results": {
    "document_type": "documenttype",
    "is_confident": true,
    "relative_confidence": 1.234567,
    "document_type_scores": [
      {
        "document_type": "documenttype",
        "score": 1.23456
      },
      {
        "document_type": "otherdocumenttype",
        "score": 33.45678
      }
    ]
  }
}`
var exampleOutputList = "[" + exampleOutputObject + "]"

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
