package tests

import (
	"bytes"
	"context"
	"errors"
	ch360mocks "github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/cmd/surf/services"
	"github.com/CloudHub360/ch360.go/cmd/surf/services/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type parallelExtractionSuite struct {
	suite.Suite
	sut              *services.ParallelExtractionService
	fileExtractor    *mocks.FileExtractor
	documentGetter   *ch360mocks.DocumentGetter
	extractorName    string
	documentId       string
	extractionResult *results.ExtractionResult
	testFilePatterns []string
	output           *bytes.Buffer
	progressHandler  *mocks.ProgressHandler
	ctx              context.Context
}

func (suite *parallelExtractionSuite) SetupTest() {
	suite.extractorName = generators.String("extractorName")
	suite.documentId = generators.String("documentId")
	suite.extractionResult = &results.ExtractionResult{}
	suite.testFilePatterns = []string{"testdata/empty-file1.txt", "testdata/empty-file2.txt"}

	suite.fileExtractor = new(mocks.FileExtractor)
	suite.documentGetter = new(ch360mocks.DocumentGetter)

	suite.output = &bytes.Buffer{}
	suite.ctx, _ = context.WithCancel(context.Background())

	suite.progressHandler = new(mocks.ProgressHandler)

	suite.documentGetter.
		On("GetAll", mock.Anything).
		Return(nil, nil)

	suite.progressHandler.
		On("NotifyStart", mock.Anything).
		Return(nil)

	suite.progressHandler.
		On("Notify", mock.Anything, mock.Anything).
		Return(nil)

	suite.progressHandler.
		On("NotifyFinish").
		Return(nil)

	suite.progressHandler.
		On("NotifyErr", mock.Anything, mock.Anything).
		Return(nil)

	suite.fileExtractor.
		On("Extract", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	suite.sut = services.NewParallelExtractionService(suite.fileExtractor, suite.documentGetter,
		suite.progressHandler)
}

func TestExtractSuiteRunner(t *testing.T) {
	suite.Run(t, new(parallelExtractionSuite))
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Processes_All_Files_Matched_By_Pattern() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.ExtractAll(suite.ctx, suite.testFilePatterns, suite.extractorName)

	suite.fileExtractor.AssertNumberOfCalls(suite.T(), "Extract", expectedCallCount)
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Calls_ProgressHandler_NotifyStart() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.ExtractAll(suite.ctx, suite.testFilePatterns, suite.extractorName)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyStart", expectedCallCount)
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Calls_ProgressHandler_Notify_For_Each_File() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.ExtractAll(suite.ctx, suite.testFilePatterns, suite.extractorName)

	suite.progressHandler.AssertNumberOfCalls(suite.T(), "Notify", expectedCallCount)
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Calls_ProgressHandler_NotifyFinish() {
	_ = suite.sut.ExtractAll(suite.ctx, suite.testFilePatterns, suite.extractorName)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyFinish")
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Return_Nil_On_Success() {
	err := suite.sut.ExtractAll(suite.ctx, suite.testFilePatterns, suite.extractorName)

	assert.Nil(suite.T(), err)
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Returns_Error_If_File_Does_Not_Exist() {
	nonExistentFile := "non-existent-file.pdf"

	err := suite.sut.ExtractAll(suite.ctx, []string{nonExistentFile}, suite.extractorName)

	assert.Error(suite.T(), err)
	suite.fileExtractor.AssertNumberOfCalls(suite.T(), "Extract", 0)
}

func (suite *parallelExtractionSuite) Test_ExtractAll_Returns_Error_If_ExtractDocument_Fails() {
	suite.fileExtractor.ExpectedCalls = nil
	extractErr := errors.New("simulated error")

	suite.fileExtractor.
		On("Extract", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, extractErr)

	err := suite.sut.ExtractAll(suite.ctx, suite.testFilePatterns, suite.extractorName)

	assert.Error(suite.T(), err)
}
