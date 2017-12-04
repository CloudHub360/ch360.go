package ch360

import (
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Client_Calls_Sender_With_Correct_Url(t *testing.T) {
	// Arrange
	sender := new(mocks.Sender)
	sender.On("Send", mock.Anything, mock.Anything, nil).Return(nil, nil)

	sut := ClassifiersClient{
		sender: sender,
	}
	classifierName := "classifier-name"

	// Act
	sut.CreateClassifier(classifierName)

	// Assert
	sender.AssertCalled(t, "Send", "POST", "/classifiers/"+classifierName, nil)
}
