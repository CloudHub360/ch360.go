package tests

import (
	"bytes"
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

type TableResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.TableClassifyResultsFormatter
	filename string
	result   *results.ClassificationResult
}

func (suite *TableResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewTableClassifyResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = &results.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestTableResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(TableResultsFormatterSuite))
}

func (suite *TableResultsFormatterSuite) TestWriteResult_With_Header_Option_Writes_Table_Header() {
	suite.sut.WriteResult(suite.output, suite.filename, suite.result, formatters.IncludeHeader)

	header := fmt.Sprintf(formatters.TableClassifyFormatterOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	assert.True(suite.T(), strings.HasPrefix(suite.output.String(), header))
}

func (suite *TableResultsFormatterSuite) TestWrites_ResultWithCorrectFormat() {
	expectedOutput := "document1.tif                        documenttype                     true\n"
	filename := "document1.tif"
	result := &results.ClassificationResult{
		DocumentType: "documenttype",
		IsConfident:  true,
	}

	err := suite.sut.WriteResult(suite.output, filename, result, 0)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedOutput, suite.output.String())
}

func (suite *TableResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *TableResultsFormatterSuite) TestWrites_DocumentType() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *TableResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), "false"))
}

func (suite *TableResultsFormatterSuite) TestWrites_Filename_Only_When_It_Has_Path() {
	var filename string
	if runtime.GOOS == "windows" { //So tests can run on both Windows & Linux
		filename = `C:\folder\document1.tif`
	} else {
		filename = `/var/something/document1.tif`
	}

	expectedFilename := `document1.tif`

	err := suite.sut.WriteResult(suite.output, filename, suite.result, 0)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedFilename, suite.output.String()[:len(expectedFilename)])
}

func (suite *TableResultsFormatterSuite) TestFlush_Writes_Nothing() {
	suite.sut.Flush(suite.output)

	assert.Equal(suite.T(), "", suite.output.String())
}
