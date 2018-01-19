package tests

import (
	"bytes"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"runtime"
	"strings"
	"testing"
)

type JsonResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.JsonClassifyResultsFormatter
	filename string
	result   *types.ClassificationResult
}

func (suite *JsonResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewJsonClassifyResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = &types.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestJsonTableResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(JsonResultsFormatterSuite))
}

func (suite *JsonResultsFormatterSuite) TestWriteHeader_Writes_Header() {
	suite.sut.WriteHeader(suite.output)

	assert.Equal(suite.T(), "[", suite.output.String())
}

func (suite *JsonResultsFormatterSuite) TestWrites_ResultWithCorrectFormat() {
	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, exampleFilename, exampleResult)
	suite.sut.WriteFooter(suite.output)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), exampleOutput, suite.output.String())
}

func (suite *JsonResultsFormatterSuite) TestWrites_Filename() {
	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *JsonResultsFormatterSuite) TestWrites_DocumentType() {
	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *JsonResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *JsonResultsFormatterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
	var filename string
	var expectedFilename string

	if runtime.GOOS == "windows" {
		filename = `C:/folder/document1.tif`
		expectedFilename = `C:\\folder\\document1.tif` //Json pretty-printing escapes slashes
	} else {
		filename = `/var/folder/document1.tif`
		expectedFilename = `/var/folder/document1.tif`
	}

	suite.sut.WriteHeader(suite.output)
	err := suite.sut.WriteResult(suite.output, filename, suite.result)

	require.Nil(suite.T(), err)
	containsFilenameWithPath := strings.Contains(suite.output.String(), expectedFilename)

	assert.True(suite.T(), containsFilenameWithPath)
	if !containsFilenameWithPath { // To aid debugging if test fails
		fmt.Println(suite.output.String())
		fmt.Printf("Output does not contain %s", expectedFilename)
	}
}

//TODO Test that WriteFooter writes nothing

var exampleFilename = "document1.tif"
var exampleResult = &types.ClassificationResult{
	DocumentType:       "documenttype",
	IsConfident:        true,
	RelativeConfidence: 1.234567,
	DocumentTypeScores: []types.DocumentTypeScore{
		types.DocumentTypeScore{
			DocumentType: "documenttype",
			Score:        1.23456,
		},
		types.DocumentTypeScore{
			DocumentType: "otherdocumenttype",
			Score:        33.45678,
		},
	}}

var exampleOutput = `[{
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
}]`
