package ch360_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FileClassifierSuite struct {
	suite.Suite
	sut                  *ch360.FileClassifier
	documentCreator      *mocks.DocumentCreator
	documentDeleter      *mocks.DocumentDeleter
	documentClassifier   *mocks.DocumentClassifier
	documentGetter       *mocks.DocumentGetter
	classifierName       string
	documentId           string
	document             ch360.Document
	classificationResult *results.ClassificationResult
	testFileContent      []byte
	testFileContentBuf   *bytes.Buffer
	ctx                  context.Context
}

func (suite *FileClassifierSuite) SetupTest() {
	suite.classifierName = generators.String("classifier-name")
	suite.documentId = generators.String("documentId")
	suite.document = ch360.Document{
		Id: suite.documentId,
	}
	suite.classificationResult = &results.ClassificationResult{}
	suite.testFileContent = generators.Bytes()
	suite.testFileContentBuf = bytes.NewBuffer(suite.testFileContent)

	suite.documentCreator = new(mocks.DocumentCreator)
	suite.documentClassifier = new(mocks.DocumentClassifier)
	suite.documentDeleter = new(mocks.DocumentDeleter)
	suite.documentGetter = new(mocks.DocumentGetter)
	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return(suite.document, nil)
	suite.documentClassifier.On("Classify", mock.Anything, mock.Anything,
		mock.Anything).Return(suite.classificationResult, nil)
	suite.documentDeleter.On("Delete", mock.Anything, mock.Anything).Return(nil)
	suite.documentGetter.On("GetAll", mock.Anything).Return(nil, nil)

	suite.ctx, _ = context.WithCancel(context.Background())

	suite.sut = ch360.NewFileClassifier(suite.documentCreator, suite.documentClassifier, suite.documentDeleter)
}

func TestFileClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(FileClassifierSuite))
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Calls_Create_Document_With_File_Content() {
	_, err := suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", mock.Anything, suite.testFileContentBuf)
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Calls_Create_Document_With_Background_Context() {
	_, err := suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", context.Background(), mock.Anything)
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Calls_Classify_With_DocumentId_And_ClassifierName() {
	_, err := suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentClassifier.AssertCalled(suite.T(), "Classify", mock.Anything, suite.documentId, suite.classifierName)
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Calls_Delete_With_DocumentId() {
	_, err := suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Calls_Delete_With_Background_Context() {
	_, err := suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", context.Background(), mock.Anything)
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Returns_Error_If_CreateDocument_Fails() {
	suite.documentCreator.ExpectedCalls = nil
	classifyErr := errors.New("simulated error")
	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return(ch360.Document{}, classifyErr)

	_, err := suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	assert.Equal(suite.T(), classifyErr, err)
}

func (suite *FileClassifierSuite) TestFileClassifier_Classify_Deletes_Document_If_ClassifyDocument_Fails() {
	expectedErr := errors.New("simulated error")
	suite.documentClassifier.ExpectedCalls = nil
	suite.documentClassifier.On("Classify", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	suite.sut.Classify(suite.ctx, suite.testFileContentBuf, suite.classifierName)

	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}
