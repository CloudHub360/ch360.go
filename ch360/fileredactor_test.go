package ch360_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/ch360/request"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"testing"
)

type FileRedactorSuite struct {
	suite.Suite
	sut                *ch360.FileRedactor
	documentCreator    *mocks.DocumentCreator
	documentDeleter    *mocks.DocumentDeleter
	documentRedactor   *mocks.DocumentRedactor
	documentExtractor  *mocks.DocumentExtractor
	extractorName      string
	documentId         string
	document           ch360.Document
	extractionResult   results.ExtractForRedactionResult
	redactionResult    io.ReadCloser
	testFileContent    []byte
	testFileContentBuf *bytes.Buffer
	ctx                context.Context
}

func (suite *FileRedactorSuite) SetupTest() {
	suite.extractorName = generators.String("extractor-name")
	suite.documentId = generators.String("documentId")
	suite.document = ch360.Document{
		Id: suite.documentId,
	}
	suite.extractionResult = results.ExtractForRedactionResult{
		Marks:      nil,
		ApplyMarks: false,
		Bookmarks:  nil,
	}
	suite.redactionResult = ioutil.NopCloser(bytes.NewBuffer(generators.Bytes()))
	suite.testFileContent = generators.Bytes()
	suite.testFileContentBuf = bytes.NewBuffer(suite.testFileContent)

	suite.documentCreator = new(mocks.DocumentCreator)
	suite.documentRedactor = new(mocks.DocumentRedactor)
	suite.documentDeleter = new(mocks.DocumentDeleter)
	suite.documentExtractor = new(mocks.DocumentExtractor)
	suite.documentCreator.
		On("Create", mock.Anything, mock.Anything).
		Return(suite.document, nil)
	suite.documentRedactor.
		On("Redact", mock.Anything, mock.Anything, mock.Anything).
		Return(suite.redactionResult, nil)
	suite.documentDeleter.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)
	suite.documentExtractor.
		On("ExtractForRedaction", mock.Anything, mock.Anything, mock.Anything).
		Return(&suite.extractionResult, nil)

	suite.ctx, _ = context.WithCancel(context.Background())

	suite.sut = ch360.NewFileRedactor(suite.documentCreator,
		suite.documentExtractor,
		suite.documentRedactor,
		suite.documentDeleter)
}

func TestFileRedactorSuiteRunner(t *testing.T) {
	suite.Run(t, new(FileRedactorSuite))
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Calls_Create_Document_With_File_Content() {
	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", mock.Anything, suite.testFileContentBuf)
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Calls_Create_Document_With_Background_Context() {
	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentCreator.AssertCalled(suite.T(), "Create", context.Background(), mock.Anything)
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Calls_Redact_With_DocumentId_And_ExtractionResult() {
	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentRedactor.AssertCalled(suite.T(), "Redact", mock.Anything, suite.documentId,
		(request.RedactedPdfRequest)(suite.extractionResult))
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Calls_Extract_With_DocumentId_And_ExtractorName() {
	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentExtractor.AssertCalled(suite.T(), "ExtractForRedaction", mock.Anything,
		suite.documentId,
		suite.extractorName)
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Calls_Delete_With_DocumentId() {
	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Calls_Delete_With_Background_Context() {
	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Nil(suite.T(), err)
	suite.documentDeleter.AssertCalled(suite.T(), "Delete", context.Background(), mock.Anything)
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Returns_Error_If_CreateDocument_Fails() {
	suite.documentCreator.ExpectedCalls = nil
	redactErr := errors.New("simulated error")
	suite.documentCreator.On("Create", mock.Anything, mock.Anything).Return(ch360.Document{}, redactErr)

	_, err := suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	assert.Equal(suite.T(), redactErr, err)
}

func (suite *FileRedactorSuite) TestFileRedactor_Redact_Deletes_Document_If_RedactDocument_Fails() {
	expectedErr := errors.New("simulated error")
	suite.documentRedactor.ExpectedCalls = nil
	suite.documentRedactor.On("Redact", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	suite.sut.Redact(suite.ctx, suite.testFileContentBuf, suite.extractorName)

	suite.documentDeleter.AssertCalled(suite.T(), "Delete", mock.Anything, suite.documentId)
}
