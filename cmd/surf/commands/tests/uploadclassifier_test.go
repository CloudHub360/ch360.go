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
	"os"
	"testing"
)

type UploadClassifierSuite struct {
	suite.Suite
	output         *bytes.Buffer
	uploader       *mocks.ClassifierUploader
	sut            *commands.UploadClassifierCmd
	classifierName string
	classifierFile *os.File
	ctx            context.Context
}

func (suite *UploadClassifierSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.uploader = new(mocks.ClassifierUploader)

	suite.classifierName = generators.String("classifier-name")
	suite.classifierFile, _ = os.Open("testdata/emptyclassifier.clf")

	suite.sut = suite.anUploadClassifierCommandWithFile(suite.classifierFile)

	suite.uploader.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.ctx = context.Background()
}

func (suite *UploadClassifierSuite) anUploadClassifierCommandWithFile(classifierFile *os.File) *commands.UploadClassifierCmd {
	return &commands.UploadClassifierCmd{
		Uploader:           suite.uploader,
		ClassifierName:     suite.classifierName,
		ClassifierContents: classifierFile,
	}
}

func TestUploadClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(UploadClassifierSuite))
}

func (suite *UploadClassifierSuite) ClearExpectedCalls() {
	suite.uploader.ExpectedCalls = nil
}

func (suite *UploadClassifierSuite) TestUploadClassifier_Execute_Uploads_The_Named_Classifier() {
	suite.sut.Execute(context.Background())

	suite.uploader.AssertCalled(suite.T(), "Upload", suite.ctx, suite.classifierName, suite.classifierFile)
}

func (suite *UploadClassifierSuite) TestUploadClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Uploaded() {
	expected := errors.New("Failed")
	suite.ClearExpectedCalls()
	suite.uploader.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(expected)

	err := suite.sut.Execute(suite.ctx)

	assert.NotNil(suite.T(), err)
	suite.uploader.AssertCalled(suite.T(), "Upload", suite.ctx, suite.classifierName, suite.classifierFile)
}
