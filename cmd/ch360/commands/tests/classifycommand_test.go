package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	cmdmocks "github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go/build"
	"io/ioutil"
	"testing"
)

type ClassifySuite struct {
	suite.Suite
	sut                  *commands.ClassifyCommand
	documentCreator      *mocks.DocumentCreator
	documentDeleter      *mocks.DocumentDeleter
	documentClassifier   *mocks.DocumentClassifier
	documentGetter       *mocks.DocumentGetter
	classifierName       string
	documentId           string
	classificationResult *types.ClassificationResult
	testFilePath         string
	testFilesPattern     string
	output               *bytes.Buffer
	resultsWriter        *cmdmocks.ClassifyResultsWriter
	ctx                  context.Context
}

func (suite *ClassifySuite) SetupTest() {
	suite.classifierName = generators.String("classifiername")
	suite.documentId = generators.String("documentId")
	suite.classificationResult = &types.ClassificationResult{
		DocumentType: generators.String("documentType"),
		IsConfident:  generators.Bool(),
	}
	suite.testFilePath = build.Default.GOPATH + "/src/github.com/CloudHub360/ch360.go/test/documents/document1.pdf"
	suite.testFilesPattern = build.Default.GOPATH + "/src/github.com/CloudHub360/ch360.go/test/documents/**/*.pdf"

	suite.documentCreator = new(mocks.DocumentCreator)
	suite.documentClassifier = new(mocks.DocumentClassifier)
	suite.documentDeleter = new(mocks.DocumentDeleter)
	suite.documentGetter = new(mocks.DocumentGetter)
	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return(suite.documentId, nil)
	suite.documentClassifier.On("Classify", mock.Anything, mock.Anything, mock.Anything).Return(suite.classificationResult, nil)
	suite.documentDeleter.On("Delete", mock.Anything, mock.Anything).Return(nil)
	suite.documentGetter.On("GetAll", mock.Anything).Return(nil, nil)

	suite.output = &bytes.Buffer{}
	suite.ctx, _ = context.WithCancel(context.Background())

	suite.resultsWriter = new(cmdmocks.ClassifyResultsWriter)
	suite.resultsWriter.On("Start").Return(nil)
	suite.resultsWriter.On("WriteResult", mock.Anything, mock.Anything).Return(nil)
	suite.resultsWriter.On("Finish").Return(nil)

	suite.sut = commands.NewClassifyCommand(
		suite.resultsWriter,
		suite.output,
		suite.documentClassifier,
		suite.documentCreator,
		suite.documentDeleter,
		suite.documentGetter,
		10)
}

func TestClassifySuiteRunner(t *testing.T) {
	suite.Run(t, new(ClassifySuite))
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Create_Document_With_File_Content() {
	expectedContents, err := ioutil.ReadFile(suite.testFilePath)
	require.Nil(suite.T(), err)

	err = suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", mock.Anything, expectedContents)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Create_Document_With_Background_Context() {
	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", context.Background(), mock.Anything)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Classify_With_DocumentId_And_ClassifierName() {
	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentClassifier.AssertCalled(suite.T(), "Classify", mock.Anything, suite.documentId, suite.classifierName)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Delete_With_DocumentId() {
	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Delete_With_Background_Context() {
	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", context.Background(), mock.Anything)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Processes_All_Files_Matched_By_Pattern() {
	suite.sut.Execute(suite.ctx, suite.testFilesPattern, suite.classifierName)

	suite.documentClassifier.AssertNumberOfCalls(suite.T(), "Classify", 5)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_ResultsWriter_Start() {
	suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	require.True(suite.T(), len(suite.resultsWriter.Calls) > 0)
	assert.Equal(suite.T(), "Start", suite.resultsWriter.Calls[0].Method)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_ResultsWriter_Write_For_Each_File() {
	suite.sut.Execute(suite.ctx, suite.testFilesPattern, suite.classifierName)

	// There are 5 files identified by suite.testFilesPattern
	suite.resultsWriter.AssertNumberOfCalls(suite.T(), "WriteResult", 5)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_ResultsWriter_Write_With_Correct_Parameters() {
	suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	resultsCall := suite.resultsWriter.Calls[1]
	assert.Equal(suite.T(), "WriteResult", resultsCall.Method)
	suite.AssertWriteResultsCallHasCorrectParameters(resultsCall)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_ResultsWriter_Finish() {
	suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	require.Equal(suite.T(), 3, len(suite.resultsWriter.Calls))
	assert.Equal(suite.T(), "Finish", suite.resultsWriter.Calls[2].Method)
}

func (suite *ClassifySuite) AssertWriteResultsCallHasCorrectParameters(call mock.Call) {
	require.Equal(suite.T(), 2, len(call.Arguments))
	assert.Equal(suite.T(), suite.testFilePath, call.Arguments[0])
	assert.Equal(suite.T(), suite.classificationResult, call.Arguments[1])
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Return_Nil_On_Success() {
	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)
	assert.Nil(suite.T(), err)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Specific_Error_If_File_Does_Not_Exist() {
	nonExistentFile := build.Default.GOPATH + "/non-existentfile.pdf"
	expectedErr := errors.New(fmt.Sprintf("File %s does not exist", nonExistentFile))

	err := suite.sut.Execute(suite.ctx, nonExistentFile, suite.classifierName)

	assert.Equal(suite.T(), expectedErr, err)
	suite.documentCreator.AssertNotCalled(suite.T(), "Create", mock.Anything, mock.Anything)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Error_If_ReadFile_Fails() {
	nonExistentFile := build.Default.GOPATH + "/non-existentfile.pdf"
	err := suite.sut.Execute(suite.ctx, nonExistentFile, suite.classifierName)

	assert.NotNil(suite.T(), err)
	suite.documentCreator.AssertNotCalled(suite.T(), "Create", mock.Anything, mock.Anything)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Error_If_CreateDocument_Fails() {
	suite.documentCreator.ExpectedCalls = nil
	classifyErr := errors.New("simulated error")
	expectedErr := errors.New(fmt.Sprintf("Error classifying file %s: %s", suite.testFilePath, classifyErr.Error()))
	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return("", classifyErr)
	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Error_If_ClassifyDocument_Fails() {
	suite.documentClassifier.ExpectedCalls = nil
	classifyErr := errors.New("simulated error")
	expectedErr := errors.New(fmt.Sprintf("Error classifying file %s: %s", suite.testFilePath, classifyErr.Error()))
	suite.documentClassifier.On("Classify", mock.Anything, mock.Anything, mock.Anything).Return(nil, classifyErr)

	err := suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *ClassifySuite) TestClassifyDoer_Deletes_Document_If_ClassifyDocument_Fails() {
	expectedErr := errors.New("simulated error")
	suite.documentClassifier.ExpectedCalls = nil
	suite.documentClassifier.On("Classify", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	suite.sut.Execute(suite.ctx, suite.testFilePath, suite.classifierName)

	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}
