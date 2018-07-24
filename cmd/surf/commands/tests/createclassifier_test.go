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

type CreateClassifierSuite struct {
	suite.Suite
	output         *bytes.Buffer
	deleter        *mocks.ClassifierDeleter
	trainer        *mocks.ClassifierTrainer
	creator        *mocks.ClassifierCreator
	sut            *commands.CreateClassifier
	classifierName string
	samplesArchive *os.File
	ctx            context.Context
}

func (suite *CreateClassifierSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.deleter = new(mocks.ClassifierDeleter)
	suite.trainer = new(mocks.ClassifierTrainer)
	suite.creator = new(mocks.ClassifierCreator)

	suite.classifierName = generators.String("classifier-name")
	suite.samplesArchive, _ = os.Open("testdata/samples.zip")

	suite.sut = suite.aClassifierCommandWithSamplesArchive(suite.samplesArchive)

	suite.creator.On("Create", mock.Anything, mock.Anything).Return(nil)
	suite.trainer.On("Train", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.deleter.On("Delete", mock.Anything, mock.Anything).Return(nil)
	suite.ctx = context.Background()
}

func (suite *CreateClassifierSuite) aClassifierCommandWithSamplesArchive(samplesArchive *os.File) *commands.CreateClassifier {
	return commands.NewCreateClassifier(suite.output,
		suite.creator,
		suite.trainer,
		suite.deleter,
		suite.classifierName,
		samplesArchive)
}

func TestCreateClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateClassifierSuite))
}

func (suite *CreateClassifierSuite) ClearExpectedCalls(creator, trainer, deleter bool) {
	if creator {
		suite.creator.ExpectedCalls = nil
	}

	if trainer {
		suite.trainer.ExpectedCalls = nil
	}

	if deleter {
		suite.deleter.ExpectedCalls = nil
	}
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Creates_The_Named_Classifier() {
	suite.sut.Execute(suite.ctx)

	suite.creator.AssertCalled(suite.T(), "Create", suite.ctx, suite.classifierName)
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Trains_The_New_Classifier() {
	suite.sut.Execute(suite.ctx)

	suite.trainer.AssertCalled(suite.T(), "Train", suite.ctx, suite.classifierName, suite.samplesArchive)
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Created() {
	expected := errors.New("Failed")
	suite.ClearExpectedCalls(true, false, false)
	suite.creator.On("Create", mock.Anything, mock.Anything).Return(expected)

	err := suite.sut.Execute(context.Background())

	assert.NotNil(suite.T(), err)
	suite.creator.AssertCalled(suite.T(), "Create", suite.ctx, suite.classifierName)
	suite.trainer.AssertNotCalled(suite.T(), "Train", suite.ctx, suite.classifierName, suite.samplesArchive)
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Returns_Error_If_The_Classifier_Cannot_Be_Created() {
	expectedErr := errors.New("Error message")
	suite.ClearExpectedCalls(true, false, false)
	suite.creator.On("Create", mock.Anything, mock.Anything).Return(expectedErr)

	receivedErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Returns_Error_If_The_Classifier_Cannot_Be_Trained() {
	expectedErr := errors.New("Error message")
	suite.ClearExpectedCalls(true, true, true)
	suite.creator.On("Create", mock.Anything, mock.Anything).Return(nil)
	suite.trainer.On("Train", mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)
	suite.deleter.On("Delete", mock.Anything, mock.Anything).Return(nil)

	receivedErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Deletes_The_Classifier_If_The_Classifier_Cannot_Be_Trained_From_The_Samples() {
	suite.deleter.On("Execute", mock.Anything).Return(nil)
	expected := errors.New("Failed")
	suite.ClearExpectedCalls(true, true, false)
	suite.creator.On("Create", mock.Anything, mock.Anything).Return(nil)
	suite.trainer.On("Train", mock.Anything, mock.Anything, mock.Anything).Return(expected)

	suite.sut.Execute(context.Background())

	suite.deleter.AssertCalled(suite.T(), "Delete", suite.ctx, suite.classifierName)
}
