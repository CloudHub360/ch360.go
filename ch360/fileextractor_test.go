package ch360_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/ch360/mocks"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/test/generators"
	"testing"
)

type FileExtractorSuite struct {
	suite.Suite
	sut                *ch360.FileExtractor
	documentCreator    *mocks.DocumentCreator
	documentDeleter    *mocks.DocumentDeleter
	documentExtractor  *mocks.DocumentExtractor
	documentGetter     *mocks.DocumentGetter
	extractorName      string
	document           ch360.Document
	documentId         string
	extractionResult   *results.ExtractionResult
	testFileContent    []byte
	testFileContentBuf *bytes.Buffer
	ctx                context.Context
}

func (suite *FileExtractorSuite) SetupTest() {
	suite.extractorName = generators.String("extractor-name")
	suite.documentId = generators.String("documentId")
	suite.document = ch360.Document{
		Id: suite.documentId,
	}
	suite.extractionResult = &results.ExtractionResult{}
	suite.testFileContent = generators.Bytes()
	suite.testFileContentBuf = bytes.NewBuffer(suite.testFileContent)

	suite.documentCreator = new(mocks.DocumentCreator)
	suite.documentExtractor = new(mocks.DocumentExtractor)
	suite.documentDeleter = new(mocks.DocumentDeleter)
	suite.documentGetter = new(mocks.DocumentGetter)

	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return(suite.document, nil)
	suite.documentExtractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(suite.extractionResult, nil)
	suite.documentDeleter.On("Delete", mock.Anything, mock.Anything).Return(nil)
	suite.documentGetter.On("GetAll", mock.Anything).Return(nil, nil)

	suite.ctx, _ = context.WithCancel(context.Background())

	suite.sut = ch360.NewFileExtractor(suite.documentCreator, suite.documentExtractor, suite.documentDeleter)
}

func TestFileExtractorSuiteRunner(t *testing.T) {
	suite.Run(t, new(FileExtractorSuite))
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Calls_Create_Document_With_File_Content() {
	_, err := suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", mock.Anything, suite.testFileContentBuf)
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Calls_Create_Document_With_Background_Context() {
	_, err := suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", context.Background(), mock.Anything)
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Calls_Extract_With_DocumentId_And_ClassifierName() {
	_, err := suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentExtractor.AssertCalled(suite.T(), "Extract", mock.Anything, suite.documentId, suite.extractorName)
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Calls_Delete_With_DocumentId() {
	_, err := suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Calls_Delete_With_Background_Context() {
	_, err := suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", context.Background(), mock.Anything)
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Returns_Error_If_CreateDocument_Fails() {
	suite.documentCreator.ExpectedCalls = nil
	extractErr := errors.New("simulated error")
	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return(ch360.Document{}, extractErr)

	_, err := suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Equal(suite.T(), extractErr, err)
}

func (suite *FileExtractorSuite) TestFileExtractor_Extract_Deletes_Document_If_ExtractDocument_Fails() {
	expectedErr := errors.New("simulated error")
	suite.documentExtractor.ExpectedCalls = nil
	suite.documentExtractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	suite.sut.Extract(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}
