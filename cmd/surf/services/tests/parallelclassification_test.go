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

type parallelClassificationSuite struct {
	suite.Suite
	sut                  *services.ParallelClassificationService
	fileClassifier       *mocks.FileClassifier
	documentGetter       *ch360mocks.DocumentGetter
	classifierName       string
	documentId           string
	classificationResult *results.ClassificationResult
	testFilePatterns     []string
	output               *bytes.Buffer
	progressHandler      *mocks.ProgressHandler
	ctx                  context.Context
}

func (suite *parallelClassificationSuite) SetupTest() {
	suite.classifierName = generators.String("classifierName")
	suite.documentId = generators.String("documentId")
	suite.classificationResult = &results.ClassificationResult{}
	suite.testFilePatterns = []string{"testdata/empty-file1.txt", "testdata/empty-file2.txt"}

	suite.fileClassifier = new(mocks.FileClassifier)
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

	suite.fileClassifier.
		On("Classify", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	suite.sut = services.NewParallelClassificationService(suite.fileClassifier, suite.documentGetter,
		suite.progressHandler)
}

func TestClassifySuiteRunner(t *testing.T) {
	suite.Run(t, new(parallelClassificationSuite))
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Processes_All_Files_Matched_By_Pattern() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.ClassifyAll(suite.ctx, suite.testFilePatterns, suite.classifierName)

	suite.fileClassifier.AssertNumberOfCalls(suite.T(), "Classify", expectedCallCount)
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Calls_ProgressHandler_NotifyStart() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.ClassifyAll(suite.ctx, suite.testFilePatterns, suite.classifierName)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyStart", expectedCallCount)
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Calls_ProgressHandler_Notify_For_Each_File() {
	var expectedCallCount = len(suite.testFilePatterns)

	_ = suite.sut.ClassifyAll(suite.ctx, suite.testFilePatterns, suite.classifierName)

	suite.progressHandler.AssertNumberOfCalls(suite.T(), "Notify", expectedCallCount)
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Calls_ProgressHandler_NotifyFinish() {
	_ = suite.sut.ClassifyAll(suite.ctx, suite.testFilePatterns, suite.classifierName)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyFinish")
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Return_Nil_On_Success() {
	err := suite.sut.ClassifyAll(suite.ctx, suite.testFilePatterns, suite.classifierName)

	assert.Nil(suite.T(), err)
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Returns_Error_If_File_Does_Not_Exist() {
	nonExistentFile := "non-existent-file.pdf"

	err := suite.sut.ClassifyAll(suite.ctx, []string{nonExistentFile}, suite.classifierName)

	assert.Error(suite.T(), err)
	suite.fileClassifier.AssertNumberOfCalls(suite.T(), "Classify", 0)
}

func (suite *parallelClassificationSuite) Test_ClassifyAll_Returns_Error_If_ClassifyDocument_Fails() {
	suite.fileClassifier.ExpectedCalls = nil
	classifyErr := errors.New("simulated error")

	suite.fileClassifier.
		On("Classify", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, classifyErr)

	err := suite.sut.ClassifyAll(suite.ctx, suite.testFilePatterns, suite.classifierName)

	assert.Error(suite.T(), err)
}
