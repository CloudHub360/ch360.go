package ch360

import (
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func Test_Client_Calls_Sender_With_Correct_Url(t *testing.T) {
	// Arrange
	sender := new(mocks.HttpDoer)
	sender.On("Do", mock.Anything).Return(nil, nil)

	sut := ClassifiersClient{
		requestSender: sender,
		baseUrl:       "baseurl",
	}
	classifierName := "classifier-name"

	// Act
	sut.CreateClassifier(classifierName)

	// Assert
	sentRequest := (sender.Calls[0].Arguments[0]).(*http.Request)
	assert.Equal(t, "POST", sentRequest.Method)
	assert.Equal(t, "baseurl/classifiers/"+classifierName, sentRequest.URL.Path)
}
