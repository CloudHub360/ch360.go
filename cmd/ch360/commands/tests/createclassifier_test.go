package tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreateClassifierSuite struct {
	suite.Suite
	output *bytes.Buffer
}

func (suite *CreateClassifierSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
}

func TestCreateClassifierSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateClassifierSuite))
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Creates_The_Named_Classifier() {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(nil)

	sut := commands.NewCreateClassifier(suite.output, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(suite.T(), "Create", "charlie")
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Trains_The_New_Classifier() {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(nil)

	sut := commands.NewCreateClassifier(suite.output, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(suite.T(), "Train", "charlie", "samples.zip")
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Created() {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Failed")
	client.On("Create", mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(suite.output, client, deleteClassifier)
	err := sut.Execute("charlie", "samples.zip")

	assert.NotNil(suite.T(), err)
	client.AssertCalled(suite.T(), "Create", "charlie")
	client.AssertNotCalled(suite.T(), "Train", "charlie", "samples.zip")
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Writes_Error_To_Output_If_The_Classifier_Cannot_Be_Created() {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Error message")
	output := &bytes.Buffer{}
	client.On("Create", mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(output, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	assert.Equal(suite.T(), fmt.Sprintf("[FAILED]\n%s\n", expected.Error()), output.String())
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Writes_Error_To_Output_If_The_Classifier_Cannot_Be_Trained() {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Error message")
	output := &bytes.Buffer{}
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(expected)
	deleteClassifier.On("Execute", mock.Anything).Return(nil)

	sut := commands.NewCreateClassifier(output, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	assert.Equal(suite.T(), fmt.Sprintf("[OK]\nAdding samples from file 'samples.zip'... [FAILED]\n%s\n", expected.Error()), output.String())
}

func (suite *CreateClassifierSuite) TestCreateClassifier_Execute_Deletes_The_Classifier_If_The_Classifier_Cannot_Be_Trained_From_The_Samples() {
	deleteClassifier := new(mocks.ClassifierCommand)
	deleteClassifier.On("Execute", mock.Anything).Return(nil)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Failed")
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(suite.output, client, deleteClassifier)
	sut.Execute("charlie", "non-existent.zip")

	deleteClassifier.AssertCalled(suite.T(), "Execute", "charlie")
}
