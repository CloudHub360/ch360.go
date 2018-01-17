package tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreateClassifierSuite struct {
	suite.Suite
	output           *bytes.Buffer
	deleteClassifier *mocks.ClassifierCommand
	client           *mocks.CreatorTrainer
}

func (suite *CreateClassifierSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.deleteClassifier = new(mocks.ClassifierCommand)
	suite.client = new(mocks.CreatorTrainer)

	suite.client.On("Create", mock.Anything).Return(nil)
	suite.client.On("Train", mock.Anything, mock.Anything).Return(nil)
}

func TestCreateClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateClassifierSuite))
}

func (suite *CreateClassifierSuite) ClearExpectedCalls() {
	suite.client.ExpectedCalls = nil
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Creates_The_Named_Classifier() {
	sut := commands.NewCreateClassifier(suite.output, suite.client, suite.deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	suite.client.AssertCalled(suite.T(), "Create", "charlie")
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Trains_The_New_Classifier() {
	sut := commands.NewCreateClassifier(suite.output, suite.client, suite.deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	suite.client.AssertCalled(suite.T(), "Train", "charlie", "samples.zip")
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Created() {
	expected := errors.New("Failed")
	suite.ClearExpectedCalls()
	suite.client.On("Create", mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(suite.output, suite.client, suite.deleteClassifier)
	err := sut.Execute("charlie", "samples.zip")

	assert.NotNil(suite.T(), err)
	suite.client.AssertCalled(suite.T(), "Create", "charlie")
	suite.client.AssertNotCalled(suite.T(), "Train", "charlie", "samples.zip")
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Writes_Error_To_Output_If_The_Classifier_Cannot_Be_Created() {
	expected := errors.New("Error message")
	suite.ClearExpectedCalls()
	suite.client.On("Create", mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(suite.output, suite.client, suite.deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	assert.Equal(suite.T(), fmt.Sprintf("[FAILED]\n%s\n", expected.Error()), suite.output.String())
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Writes_Error_To_Output_If_The_Classifier_Cannot_Be_Trained() {
	expected := errors.New("Error message")
	suite.ClearExpectedCalls()
	suite.client.On("Create", mock.Anything).Return(nil)
	suite.client.On("Train", mock.Anything, mock.Anything).Return(expected)
	suite.deleteClassifier.On("Execute", mock.Anything).Return(nil)

	sut := commands.NewCreateClassifier(suite.output, suite.client, suite.deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	assert.Equal(suite.T(), fmt.Sprintf("[OK]\nAdding samples from file 'samples.zip'... [FAILED]\n%s\n", expected.Error()), suite.output.String())
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Deletes_The_Classifier_If_The_Classifier_Cannot_Be_Trained_From_The_Samples() {
	suite.deleteClassifier.On("Execute", mock.Anything).Return(nil)
	expected := errors.New("Failed")
	suite.ClearExpectedCalls()
	suite.client.On("Create", mock.Anything).Return(nil)
	suite.client.On("Train", mock.Anything, mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(suite.output, suite.client, suite.deleteClassifier)
	sut.Execute("charlie", "non-existent.zip")

	suite.deleteClassifier.AssertCalled(suite.T(), "Execute", "charlie")
}
