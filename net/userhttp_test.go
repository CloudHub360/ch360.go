package net

import (
	"errors"
	"github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestUserAgentHttpClient_Sets_Correct_Header(t *testing.T) {
	// Arrange
	wrapped := mocks.HttpDoer{}
	wrapped.On("Do", mock.Anything).Return(nil, nil)
	sut := NewUserAgentHttpClient(&wrapped, "test-agent")
	request, _ := http.NewRequest("GET", "https://api.waives.io", nil)

	// Act
	_, _ = sut.Do(request)

	// Assert
	assert.Equal(t, "test-agent", request.Header.Get("User-Agent"))
}

func TestUserAgentHttpClient_Returns_Values_From_Wrapped_HttpDoer(t *testing.T) {
	// Arrange
	wrapped := mocks.HttpDoer{}
	expectedResponse := &http.Response{}
	expectedErr := errors.New("test error")
	wrapped.On("Do", mock.Anything).Return(expectedResponse, expectedErr)
	sut := NewUserAgentHttpClient(&wrapped, "test-agent")
	request, _ := http.NewRequest("GET", "https://api.waives.io", nil)

	// Act
	actualResponse, actualErr := sut.Do(request)

	// Assert
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedErr, actualErr)
}
