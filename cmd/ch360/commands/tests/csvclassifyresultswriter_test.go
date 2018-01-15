package tests

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type CSVResultsWriterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *commands.CSVClassifyResultsWriter
	filename string
	result   *types.ClassificationResult
}

func (suite *CSVResultsWriterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = commands.NewCSVClassifyResultsWriter(suite.output)

	suite.filename = generators.String("filename")
	suite.result = &types.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestCSVResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(CSVResultsWriterSuite))
}

func (suite *CSVResultsWriterSuite) TestStart_Does_Not_Write_Anything() {
	suite.sut.StartWriting()

	assert.Equal(suite.T(), "", suite.output.String())
}

func (suite *CSVResultsWriterSuite) TestWrite_Returns_Error_If_Start_Not_Called() {
	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	assert.NotNil(suite.T(), err)
}

func (suite *CSVResultsWriterSuite) TestWrites_ResultWithCorrectFormat() {
	filename := "document1.tif"
	result := &types.ClassificationResult{
		DocumentType: "documenttype",
		IsConfident:  true,
	}
	expectedOutput := "document1.tif,documenttype,true\n"

	suite.sut.StartWriting()
	err := suite.sut.WriteDocumentResults(filename, result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedOutput, suite.output.String())
}

func (suite *CSVResultsWriterSuite) TestWrites_Filename() {
	suite.sut.StartWriting()
	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *CSVResultsWriterSuite) TestWrites_DocumentType() {
	suite.sut.StartWriting()
	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *CSVResultsWriterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	suite.sut.StartWriting()
	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *CSVResultsWriterSuite) TestWrites_Filename_With_Path_When_It_Has_Path() {
	filename := `C:\folder\document1.tif`

	suite.sut.StartWriting()
	err := suite.sut.WriteDocumentResults(filename, suite.result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), filename, suite.output.String()[:len(filename)])
}
