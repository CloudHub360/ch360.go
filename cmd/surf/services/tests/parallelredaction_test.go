package tests

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	ch360mocks "github.com/waives/surf/ch360/mocks"
	"github.com/waives/surf/cmd/surf/services"
	"github.com/waives/surf/cmd/surf/services/mocks"
	"github.com/waives/surf/test/generators"
	"io"
	"io/ioutil"
	"testing"
)

type parallelRedactionSuite struct {
	suite.Suite
	sut              *services.ParallelRedactionService
	fileRedactor     *mocks.FileRedactor
	documentGetter   *ch360mocks.DocumentGetter
	redactorName     string
	documentId       string
	redactionResult  io.ReadCloser
	testFilePatterns []string
	output           *bytes.Buffer
	progressHandler  *mocks.ProgressHandler
	ctx              context.Context
}

func (suite *parallelRedactionSuite) SetupTest() {
	suite.redactorName = generators.String("redactorName")
	suite.documentId = generators.String("documentId")
	suite.redactionResult = ioutil.NopCloser(bytes.NewBuffer(generators.Bytes()))
	suite.testFilePatterns = []string{"testdata/empty-file1.txt", "testdata/empty-file2.txt"}

	suite.fileRedactor = new(mocks.FileRedactor)
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

	suite.fileRedactor.
		On("Redact", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	suite.sut = services.NewParallelRedactionService(suite.fileRedactor, suite.documentGetter,
		suite.progressHandler)
}

func TestRedactSuiteRunner(t *testing.T) {
	suite.Run(t, new(parallelRedactionSuite))
}

func (suite *parallelRedactionSuite) Test_Redact_Processes_All_Files_Matched_By_Pattern() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.RedactAllWithExtractor(suite.ctx, suite.testFilePatterns, suite.redactorName)

	suite.fileRedactor.AssertNumberOfCalls(suite.T(), "Redact", expectedCallCount)
}

func (suite *parallelRedactionSuite) Test_Redact_Calls_ProgressHandler_NotifyStart() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.RedactAllWithExtractor(suite.ctx, suite.testFilePatterns, suite.redactorName)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyStart", expectedCallCount)
}

func (suite *parallelRedactionSuite) Test_Redact_Calls_ProgressHandler_Notify_For_Each_File() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.RedactAllWithExtractor(suite.ctx, suite.testFilePatterns, suite.redactorName)

	suite.progressHandler.AssertNumberOfCalls(suite.T(), "Notify", expectedCallCount)
}

func (suite *parallelRedactionSuite) Test_Redact_Calls_ProgressHandler_NotifyFinish() {
	_ = suite.sut.RedactAllWithExtractor(suite.ctx, suite.testFilePatterns, suite.redactorName)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyFinish")
}

func (suite *parallelRedactionSuite) Test_Redact_Return_Nil_On_Success() {
	err := suite.sut.RedactAllWithExtractor(suite.ctx, suite.testFilePatterns, suite.redactorName)

	assert.Nil(suite.T(), err)
}

func (suite *parallelRedactionSuite) Test_Redact_Returns_Error_If_File_Does_Not_Exist() {
	nonExistentFile := "non-existent-file.pdf"

	err := suite.sut.RedactAllWithExtractor(suite.ctx, []string{nonExistentFile}, suite.redactorName)

	assert.Error(suite.T(), err)
	suite.fileRedactor.AssertNumberOfCalls(suite.T(), "Redact", 0)
}

func (suite *parallelRedactionSuite) Test_Redact_Returns_Error_If_RedactDocument_Fails() {
	suite.fileRedactor.ExpectedCalls = nil
	redactErr := errors.New("simulated error")

	suite.fileRedactor.
		On("Redact", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, redactErr)

	err := suite.sut.RedactAllWithExtractor(suite.ctx, suite.testFilePatterns, suite.redactorName)

	assert.Error(suite.T(), err)
}
