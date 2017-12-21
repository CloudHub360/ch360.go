package tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go/build"
	"io/ioutil"
	"testing"
)

type ClassifySuite struct {
	suite.Suite
	sut            *commands.ClassifyDoer
	client         *mocks.DocumentCreatorDeleterClassifier
	classifierName string
	documentId     string
	documentType   string
	testFilePath   string
	output         *bytes.Buffer
}

func (suite *ClassifySuite) SetupTest() {
	suite.classifierName = generators.String("classifiername")
	suite.documentId = generators.String("documentId")
	suite.documentType = generators.String("documentType")
	suite.testFilePath = build.Default.GOPATH + "/src/github.com/CloudHub360/ch360.go/test/documents/document1.pdf"

	suite.client = new(mocks.DocumentCreatorDeleterClassifier)
	suite.client.On("Create", mock.Anything).Return(suite.documentId, nil)
	suite.client.On("Classify", mock.Anything, mock.Anything).Return(suite.documentType, nil)
	suite.client.On("Delete", mock.Anything).Return(nil)

	suite.output = &bytes.Buffer{}
	suite.sut = commands.NewClassifyDoer(suite.output, suite.client)
}

func TestClassifySuiteRunner(t *testing.T) {
	suite.Run(t, new(ClassifySuite))
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Create_Document_With_File_Content() {
	expectedContents, err := ioutil.ReadFile(suite.testFilePath)
	assert.Nil(suite.T(), err)

	err = suite.sut.Execute(suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.client.AssertCalled(suite.T(), "Create", expectedContents)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Classify_With_DocumentId_And_ClassifierName() {
	err := suite.sut.Execute(suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.client.AssertCalled(suite.T(), "Classify", suite.documentId, suite.classifierName)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Calls_Delete_With_DocumentId() {
	err := suite.sut.Execute(suite.testFilePath, suite.classifierName)

	assert.Nil(suite.T(), err)
	suite.client.AssertCalled(suite.T(), "Delete", suite.documentId)
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Writes_DocumentType_To_StdOut() {
	suite.sut.Execute(suite.testFilePath, suite.classifierName)

	assert.Equal(suite.T(), suite.documentType+"\n", suite.output.String())
}

func (suite *ClassifySuite) TestClassifyDoer_Execute_Return_Nil_On_Success() {
	err := suite.sut.Execute(suite.testFilePath, suite.classifierName)
	assert.Nil(suite.T(), err)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Specific_Error_If_File_Does_Not_Exist() {
	nonExistentFile := build.Default.GOPATH + "/non-existentfile.pdf"
	expectedErr := errors.New(fmt.Sprintf("File %s does not exist", nonExistentFile))

	err := suite.sut.Execute(nonExistentFile, suite.classifierName)

	assert.Equal(suite.T(), expectedErr, err)
	suite.client.AssertNotCalled(suite.T(), "Create", mock.Anything)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Error_If_ReadFile_Fails() {
	nonExistentFile := build.Default.GOPATH + "/non-existentfile.pdf"
	err := suite.sut.Execute(nonExistentFile, suite.classifierName)

	assert.NotNil(suite.T(), err)
	suite.client.AssertNotCalled(suite.T(), "Create", mock.Anything)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Error_If_CreateDocument_Fails() {
	suite.client.ExpectedCalls = nil
	expectedErr := errors.New("simulated error")
	suite.client.On("Create", mock.Anything).Return("", expectedErr)
	err := suite.sut.Execute(suite.testFilePath, suite.classifierName)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *ClassifySuite) TestClassifyDoer_Returns_Error_If_ClassifyDocument_Fails() {
	suite.client.ExpectedCalls = nil
	expectedErr := errors.New("simulated error")
	suite.client.On("Create", mock.Anything).Return(suite.documentId, nil)
	suite.client.On("Classify", mock.Anything, mock.Anything).Return("", expectedErr)
	suite.client.On("Delete", mock.Anything).Return(nil)

	err := suite.sut.Execute(suite.testFilePath, suite.classifierName)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *ClassifySuite) TestClassifyDoer_Deletes_Document_If_ClassifyDocument_Fails() {
	suite.client.ExpectedCalls = nil
	expectedErr := errors.New("simulated error")
	suite.client.On("Create", mock.Anything).Return(suite.documentId, nil)
	suite.client.On("Classify", mock.Anything, mock.Anything).Return("", expectedErr)
	suite.client.On("Delete", mock.Anything).Return(nil)

	suite.sut.Execute(suite.testFilePath, suite.classifierName)

	suite.client.AssertCalled(suite.T(), "Delete", suite.documentId)
}
