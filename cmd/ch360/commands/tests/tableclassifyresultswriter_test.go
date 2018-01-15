package tests

import (
	"bytes"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type TableResultsWriterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *commands.TableClassifyResultsWriter
	filename string
	result   *types.ClassificationResult
}

func (suite *TableResultsWriterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = commands.NewTableClassifyResultsWriter(suite.output)

	suite.filename = generators.String("filename")
	suite.result = &types.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestTableResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(TableResultsWriterSuite))
}

func (suite *TableResultsWriterSuite) TestStart_Writes_Table_Header() {
	suite.sut.StartWriting()

	header := fmt.Sprintf(commands.ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	assert.Equal(suite.T(), header, suite.output.String())
}

func (suite *TableResultsWriterSuite) TestWrites_ResultWithCorrectFormat() {
	expectedOutput := "document1.tif                        documenttype                     true\n"
	filename := "document1.tif"
	result := &types.ClassificationResult{
		DocumentType: "documenttype",
		IsConfident:  true,
	}

	err := suite.sut.WriteDocumentResults(filename, result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedOutput, suite.output.String())
}

func (suite *TableResultsWriterSuite) TestWrites_Filename() {
	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *TableResultsWriterSuite) TestWrites_DocumentType() {
	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *TableResultsWriterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	err := suite.sut.WriteDocumentResults(suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *TableResultsWriterSuite) TestWrites_Filename_Only_When_It_Has_Path() {
	filename := `C:\folder\document1.tif`
	expectedFilename := `document1.tif`

	err := suite.sut.WriteDocumentResults(filename, suite.result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedFilename, suite.output.String()[:len(expectedFilename)])
}
