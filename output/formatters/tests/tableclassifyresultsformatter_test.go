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

type TableResultsFormatterSuite struct {
	suite.Suite
	output   *bytes.Buffer
	sut      *formatters.TableClassifyResultsFormatter
	filename string
	result   *types.ClassificationResult
}

func (suite *TableResultsFormatterSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.sut = formatters.NewTableClassifyResultsFormatter()

	suite.filename = generators.String("filename")
	suite.result = &types.ClassificationResult{
		DocumentType: generators.String("documenttype"),
		IsConfident:  false,
	}
}

func TestTableResultsWriterRunner(t *testing.T) {
	suite.Run(t, new(TableResultsFormatterSuite))
}

func (suite *TableResultsFormatterSuite) TestWriteHeader_Writes_Table_Header() {
	suite.sut.WriteHeader(suite.output)

	header := fmt.Sprintf(formatters.TableFormatterOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	assert.Equal(suite.T(), header, suite.output.String())
}

func (suite *TableResultsFormatterSuite) TestWrites_ResultWithCorrectFormat() {
	expectedOutput := "document1.tif                        documenttype                     true\n"
	filename := "document1.tif"
	result := &types.ClassificationResult{
		DocumentType: "documenttype",
		IsConfident:  true,
	}

	err := suite.sut.WriteResult(suite.output, filename, result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedOutput, suite.output.String())
}

func (suite *TableResultsFormatterSuite) TestWrites_Filename() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.filename))
}

func (suite *TableResultsFormatterSuite) TestWrites_DocumentType() {
	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

	require.Nil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(suite.output.String(), suite.result.DocumentType))
}

func (suite *TableResultsFormatterSuite) TestWrites_False_For_Not_IsConfident() {
	suite.result.IsConfident = false

	err := suite.sut.WriteResult(suite.output, suite.filename, suite.result)

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

	err := suite.sut.WriteResult(suite.output, filename, suite.result)

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedFilename, suite.output.String()[:len(expectedFilename)])
}

func (suite *TableResultsFormatterSuite) TestWriteFooter_Writes_Nothing() {
	suite.sut.WriteFooter(suite.output)

	assert.Equal(suite.T(), "", suite.output.String())
}
