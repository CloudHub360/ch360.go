package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	cmdmocks "github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go/build"
	"testing"
)

type ExtractFilesSuite struct {
	suite.Suite
	sut              *commands.Extract
	fileExtractor    *cmdmocks.FileExtractor
	documentGetter   *mocks.DocumentGetter
	classifierName   string
	documentId       string
	extractionResult *results.ExtractionResult
	testFilePath     string
	testFilesPattern string
	output           *bytes.Buffer
	progressHandler  *cmdmocks.ProgressHandler
	ctx              context.Context
}

func (suite *ExtractFilesSuite) SetupTest() {
	suite.classifierName = generators.String("classifiername")
	suite.documentId = generators.String("documentId")
	suite.extractionResult = &results.ExtractionResult{}
	suite.testFilePath = build.Default.GOPATH + "/src/github.com/CloudHub360/ch360.go/test/documents/extraction/document1.pdf"
	suite.testFilesPattern = build.Default.GOPATH + "/src/github.com/CloudHub360/ch360.go/test/documents/extraction/**/*.pdf"

	suite.fileExtractor = new(cmdmocks.FileExtractor)
	suite.documentGetter = new(mocks.DocumentGetter)
	suite.fileExtractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(suite.extractionResult, nil)
	suite.documentGetter.On("GetAll", mock.Anything).Return(nil, nil)

	suite.output = &bytes.Buffer{}
	suite.ctx, _ = context.WithCancel(context.Background())

	suite.progressHandler = new(cmdmocks.ProgressHandler)

	suite.progressHandler.On("NotifyStart", mock.Anything).Return(nil)
	suite.progressHandler.On("Notify", mock.Anything, mock.Anything).Return(nil)
	suite.progressHandler.On("NotifyErr", mock.Anything, mock.Anything).Return(nil)
	suite.progressHandler.On("NotifyFinish").Return(nil)

	suite.sut = suite.aExtractCmd()
}

func (suite *ExtractFilesSuite) aExtractCmd() *commands.Extract {
	return suite.aExtractCmdWith(suite.testFilePath)

}

func (suite *ExtractFilesSuite) aExtractCmdWith(filePattern string) *commands.Extract {
	return commands.NewExtractCommand(
		suite.progressHandler,
		suite.fileExtractor,
		suite.documentGetter,
		10,
		filePattern,
		suite.classifierName)
}

func (suite *ExtractFilesSuite) aExtractCmdWithFilePattern() *commands.Extract {
	return suite.aExtractCmdWith(suite.testFilesPattern)
}

func TestExtractSuiteRunner(t *testing.T) {
	suite.Run(t, new(ExtractFilesSuite))
}

func (suite *ExtractFilesSuite) TestExtract_Command_Execute_Processes_All_Files_Matched_By_Pattern() {
	suite.sut = suite.aExtractCmdWithFilePattern()
	suite.sut.Execute(suite.ctx)

	suite.fileExtractor.AssertNumberOfCalls(suite.T(), "Extract", 5)
}

func (suite *ExtractFilesSuite) TestExtract_Command_Execute_Calls_ProgressHandler_NotifyStart() {
	suite.sut.Execute(suite.ctx)

	require.True(suite.T(), len(suite.progressHandler.Calls) > 0)
	assert.Equal(suite.T(), "NotifyStart", suite.progressHandler.Calls[0].Method)
}

func (suite *ExtractFilesSuite) TestExtract_Command_Execute_Calls_ProgressHandler_Write_For_Each_File() {
	suite.sut = suite.aExtractCmdWithFilePattern()
	suite.sut.Execute(suite.ctx)

	// There are 5 files identified by suite.testFilesPattern
	suite.progressHandler.AssertNumberOfCalls(suite.T(), "Notify", 5)
}

func (suite *ExtractFilesSuite) TestExtract_Command_Execute_Calls_ProgressHandler_Write_With_Correct_Parameters() {
	suite.sut.Execute(suite.ctx)

	resultsCall := suite.progressHandler.Calls[1]
	assert.Equal(suite.T(), "Notify", resultsCall.Method)
	suite.AssertNotifyCallHasCorrectParameters(resultsCall)
}

func (suite *ExtractFilesSuite) TestExtract_Command_Execute_Calls_ProgressHandler_Finish() {
	suite.sut.Execute(suite.ctx)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyFinish")
}

func (suite *ExtractFilesSuite) AssertNotifyCallHasCorrectParameters(call mock.Call) {
	require.Equal(suite.T(), 2, len(call.Arguments))
	assert.Equal(suite.T(), suite.testFilePath, call.Arguments[0])
	assert.Equal(suite.T(), suite.extractionResult, call.Arguments[1])
}

func (suite *ExtractFilesSuite) TestExtract_Command_Execute_Return_Nil_On_Success() {
	err := suite.sut.Execute(suite.ctx)
	assert.Nil(suite.T(), err)
}

func (suite *ExtractFilesSuite) TestExtract_Command_Returns_Specific_Error_If_File_Does_Not_Exist() {
	nonExistentFile := build.Default.GOPATH + "/non-existentfile.pdf"

	expectedErr := errors.New(fmt.Sprintf("File %s does not exist", nonExistentFile))

	suite.sut = suite.aExtractCmdWith(nonExistentFile)
	err := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expectedErr, err)
	suite.fileExtractor.AssertNumberOfCalls(suite.T(), "Extract", 0)
}

func (suite *ExtractFilesSuite) TestExtract_Command_Returns_Error_If_ReadFile_Fails() {
	nonExistentFile := build.Default.GOPATH + "/non-existentfile.pdf"
	suite.sut = suite.aExtractCmdWith(nonExistentFile)

	err := suite.sut.Execute(suite.ctx)

	assert.NotNil(suite.T(), err)
	suite.fileExtractor.AssertNumberOfCalls(suite.T(), "Extract", 0)

}

func (suite *ExtractFilesSuite) TestExtract_Command_Returns_Error_If_ExtractDocument_Fails() {
	suite.fileExtractor.ExpectedCalls = nil
	extractErr := errors.New("simulated error")
	expectedErr := errors.New(fmt.Sprintf("Error extracting file %s: %s", suite.testFilePath, extractErr.Error()))
	suite.fileExtractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(nil, extractErr)

	err := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expectedErr, err)
}
