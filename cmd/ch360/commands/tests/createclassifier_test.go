package tests

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestCreateClassifier_Execute_Creates_The_Named_Classifier(t *testing.T) {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(nil)

	sut := commands.NewCreateClassifier(os.Stdout, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(t, "Create", "charlie")
}

func TestCreateClassifier_Execute_Trains_The_New_Classifier(t *testing.T) {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(nil)

	sut := commands.NewCreateClassifier(os.Stdout, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(t, "Train", "charlie", "samples.zip")
}

func TestCreateClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Created(t *testing.T) {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Failed")
	client.On("Create", mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(os.Stdout, client, deleteClassifier)
	err := sut.Execute("charlie", "samples.zip")

	assert.NotNil(t, err)
	client.AssertCalled(t, "Create", "charlie")
	client.AssertNotCalled(t, "Train", "charlie", "samples.zip")
}

func TestCreateClassifier_Execute_Writes_Error_To_Output_If_The_Classifier_Cannot_Be_Created(t *testing.T) {
	deleteClassifier := new(mocks.ClassifierCommand)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Error message")
	output := &bytes.Buffer{}
	client.On("Create", mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(output, client, deleteClassifier)
	sut.Execute("charlie", "samples.zip")

	assert.Equal(t, fmt.Sprintf("[FAILED]\n%s\n", expected.Error()), output.String())
}

func TestCreateClassifier_Execute_Deletes_The_Classifier_If_The_Classifier_Cannot_Be_Trained_From_The_Samples(t *testing.T) {
	deleteClassifier := new(mocks.ClassifierCommand)
	deleteClassifier.On("Execute", mock.Anything).Return(nil)
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Failed")
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(expected)

	sut := commands.NewCreateClassifier(os.Stdout, client, deleteClassifier)
	sut.Execute("charlie", "non-existent.zip")

	deleteClassifier.AssertCalled(t, "Execute", "charlie")
}
