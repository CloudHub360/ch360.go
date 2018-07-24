package tests

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"testing"
)

type UploadClassifierSuite struct {
	suite.Suite
	output         *bytes.Buffer
	uploader       *mocks.ClassifierUploader
	sut            *commands.UploadClassifier
	classifierName string
	classifierFile io.ReadCloser
}

func (suite *UploadClassifierSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.uploader = new(mocks.ClassifierUploader)

	suite.classifierName = generators.String("classifier-name")
	suite.classifierFile, _ = os.Open("classifier.clf")

	suite.sut = suite.anUploadClassifierCommandWithFile(suite.classifierFile)

	suite.uploader.On("Upload", mock.Anything, mock.Anything).Return(nil)
}

func (suite *UploadClassifierSuite) anUploadClassifierCommandWithFile(classifierFile io.ReadCloser) *commands.UploadClassifier {
	return commands.NewUploadClassifier(suite.output,
		suite.uploader,
		suite.classifierName,
		classifierFile)
}

func TestUploadClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(UploadClassifierSuite))
}

func (suite *UploadClassifierSuite) ClearExpectedCalls() {
	suite.uploader.ExpectedCalls = nil
}

func (suite *UploadClassifierSuite) TestUploadClassifier_Execute_Uploads_The_Named_Classifier() {
	suite.sut.Execute(context.Background())

	suite.uploader.AssertCalled(suite.T(), "Upload", suite.classifierName, suite.classifierFile)
}

func (suite *UploadClassifierSuite) TestUploadClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Uploaded() {
	expected := errors.New("Failed")
	suite.ClearExpectedCalls()
	suite.uploader.On("Upload", mock.Anything, mock.Anything).Return(expected)

	err := suite.sut.Execute(context.Background())

	assert.NotNil(suite.T(), err)
	suite.uploader.AssertCalled(suite.T(), "Upload", suite.classifierName, suite.classifierFile)
}
