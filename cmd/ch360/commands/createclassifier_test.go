package commands

import (
	"errors"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateClassifier_Execute_Creates_The_Named_Classifier(t *testing.T) {
	client := new(mocks.Creator)
	client.On("Create", mock.Anything).Return(nil)

	sut := NewCreateClassifier(client)
	sut.Execute("charlie")

	client.AssertCalled(t, "Create", "charlie")
}

func TestCreateClassifier_Execute_Returns_An_Error_If_The_Classifier_Cannot_Be_Created(t *testing.T) {
	client := new(mocks.Creator)
	expected := errors.New("Failed")
	client.On("Create", mock.Anything).Return(expected)

	sut := NewCreateClassifier(client)
	actual := sut.Execute("charlie")

	assert.Equal(t, expected, actual)
}
