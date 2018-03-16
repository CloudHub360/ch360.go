package ch360_test

import (
	"bytes"
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"testing"
)

type fileReaderSuite struct {
	suite.Suite
	sut          *ch360.FileReader
	docCreator   *mocks.DocumentCreator
	docReader    *mocks.DocumentReader
	docDeleter   *mocks.DocumentDeleter
	fileContents []byte
	reader       io.Reader
	ctx          context.Context
	documentId   string
}

func (suite *fileReaderSuite) SetupTest() {
	suite.docCreator = &mocks.DocumentCreator{}
	suite.docReader = &mocks.DocumentReader{}
	suite.docDeleter = &mocks.DocumentDeleter{}

	suite.sut = ch360.NewFileReader(suite.docCreator, suite.docReader, suite.docDeleter)

	suite.fileContents = generators.Bytes()
	suite.reader = bytes.NewBuffer(suite.fileContents)

	suite.documentId = generators.String("documentId")

	suite.ctx, _ = context.WithCancel(context.Background())

	// set up the happy path
	suite.docCreator.
		On("Create", mock.Anything, mock.Anything).
		Return(suite.documentId, nil)
	suite.docReader.
		On("Read", mock.Anything, mock.Anything).
		Return(nil)
	suite.docReader.
		On("ReadResult", mock.Anything, mock.Anything, mock.Anything).
		Return(ioutil.NopCloser(suite.reader), nil)
	suite.docDeleter.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)
}

func TestFileReaderSuiteRunner(t *testing.T) {
	suite.Run(t, new(fileReaderSuite))
}

func (suite *fileReaderSuite) Test_DocCreator_Called_With_File_Content() {
	suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	suite.docCreator.
		AssertCalled(suite.T(), "Create", context.Background(), suite.fileContents)
}

func (suite *fileReaderSuite) Test_Returns_Error_From_DocCreator() {
	expectedErr := errors.New("generated err")
	suite.docCreator.ExpectedCalls = nil
	suite.docCreator.
		On("Create", mock.Anything, mock.Anything).
		Return("", expectedErr)

	_, receivedErr := suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	suite.Assert().Equal(expectedErr, receivedErr)
}

func (suite *fileReaderSuite) Test_Read_And_ReadResult_Called_With_Correct_Params() {
	suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	suite.docReader.
		AssertCalled(suite.T(), "Read", suite.ctx, suite.documentId)
	suite.docReader.
		AssertCalled(suite.T(), "ReadResult", suite.ctx, suite.documentId, ch360.ReadPDF)
}

func (suite *fileReaderSuite) Test_ReadResult_Not_Called_If_Read_Returns_Err() {
	// Arrange
	expectedErr := errors.New("generated err")
	suite.docReader.ExpectedCalls = nil
	suite.docReader.
		On("Read", mock.Anything, mock.Anything).
		Return(expectedErr)

	// Act
	_, receivedErr := suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	// Assert
	suite.Assert().Equal(expectedErr, receivedErr)
	suite.docReader.AssertCalled(suite.T(), "Read", suite.ctx, suite.documentId)
	suite.docReader.AssertNumberOfCalls(suite.T(), "ReadResult", 0)
}

func (suite *fileReaderSuite) Test_Delete_Called() {
	// Act
	suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	// Assert
	suite.docDeleter.AssertCalled(suite.T(), "Delete", context.Background(), suite.documentId)
}

func (suite *fileReaderSuite) Test_Delete_Called_When_Read_Returns_Error() {
	// Arrange
	expectedErr := errors.New("generated err")
	suite.docReader.ExpectedCalls = nil
	suite.docReader.
		On("Read", mock.Anything, mock.Anything).
		Return(expectedErr)

	// Act
	suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	// Assert
	suite.docDeleter.AssertCalled(suite.T(), "Delete", context.Background(), suite.documentId)
}

func (suite *fileReaderSuite) Test_Delete_Not_Called_If_Create_Returns_Error() {
	// Arrange
	expectedErr := errors.New("generated err")
	suite.docCreator.ExpectedCalls = nil
	suite.docCreator.
		On("Create", mock.Anything, mock.Anything).
		Return("", expectedErr)

	// Act
	suite.sut.Read(suite.ctx, suite.reader, ch360.ReadPDF)

	// Assert
	suite.docDeleter.AssertNumberOfCalls(suite.T(), "Delete", 0)
}
