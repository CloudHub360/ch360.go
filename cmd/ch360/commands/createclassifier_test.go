package commands

import (
	"errors"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateClassifier_Execute_Creates_The_Named_Classifier(t *testing.T) {
	client := new(mocks.CreatorTrainer)
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(nil)

	sut := NewCreateClassifier(client)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(t, "Create", "charlie")
}

func TestCreateClassifier_Execute_Trains_The_New_Classifier(t *testing.T) {
	client := new(mocks.CreatorTrainer)
	client.On("Create", mock.Anything).Return(nil)
	client.On("Train", mock.Anything, mock.Anything).Return(nil)

	sut := NewCreateClassifier(client)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(t, "Train", "charlie", "samples.zip")
}

func TestCreateClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Created(t *testing.T) {
	client := new(mocks.CreatorTrainer)
	expected := errors.New("Failed")
	client.On("Create", mock.Anything).Return(expected)

	sut := NewCreateClassifier(client)
	sut.Execute("charlie", "samples.zip")

	client.AssertCalled(t, "Create", "charlie")
	client.AssertNotCalled(t, "Train", "charlie", "samples.zip")
}
