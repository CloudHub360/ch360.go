//go:generate mockery -name FormPoster -recursive

package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"bytes"
	"github.com/CloudHub360/ch360.go/response"
	"errors"
)

var fakeClientId = "fake-client-id"
var fakeClientSecret = "fake-client-secret"

func Test_HttpTokenRetriever_Sends_Client_Id_And_Secret(t *testing.T) {
	// Arrange
	var receivedClientId string
	var receivedClientSecret string

	requestHandler := func(w http.ResponseWriter, r *http.Request) {
		receivedClientId = r.FormValue("client_id")
		receivedClientSecret = r.FormValue("client_secret")
		w.WriteHeader(200)
	}

	// create test server with requestHandler
	server := httptest.NewServer(http.HandlerFunc(requestHandler))
	defer server.Close()

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, server.Client(), server.URL, &response.ErrorChecker{})

	// Act
	sut.RetrieveToken()

	// Assert
	assert.Equal(t, fakeClientId, receivedClientId)
	assert.Equal(t, fakeClientSecret, receivedClientSecret)
}

func Test_HttpTokenRetriever_Returns_Error_On_HttpClient_Error(t *testing.T) {
	// Arrange
	tokenGetter := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, &http.Client{}, "http://invalid-url:-1", &response.ErrorChecker{})

	// Act
	_, err := tokenGetter.RetrieveToken()

	// Assert
	assert.NotNil(t, err)
}

func AnHttpResponse(body []byte, status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
	}
}

func Test_HttpTokenRetriever_Passes_Response_To_Checker(t *testing.T) {
	// Arrange
	expectedResponseBody := []byte(`{"access_token": "tokenvalue"}`)

	response := AnHttpResponse(expectedResponseBody, 200)

	mockHttpClient := new(mocks.FormPoster)
	mockResponseChecker := new(mocks.Checker)

	mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return(expectedResponseBody, nil)

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, mockHttpClient, "notused", mockResponseChecker)

	// Act
	sut.RetrieveToken()

	// Assert
	mockResponseChecker.AssertCalled(t, "Check", response, 200)
}

func Test_HttpTokenRetriever_Returns_Error_On_ResponseChecker_Error(t *testing.T) {
	// Arrange
	response := AnHttpResponse(nil, 200)

	mockHttpClient := new(mocks.FormPoster)
	mockResponseChecker := new(mocks.Checker)

	mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return(nil, errors.New("An error"))

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, mockHttpClient, "notused", mockResponseChecker)

	// Act
	_, err := sut.RetrieveToken()

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "An error")
}

func Test_HttpTokenRetriever_Parses_Token_Response(t *testing.T) {
	// Arrange
	expectedToken := "fake-token"
	fakeServer := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token": "%s"}`, expectedToken)
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(fakeServer))
	defer ts.Close()

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, ts.Client(), ts.URL, &response.ErrorChecker{})

	// Act
	token, err := sut.RetrieveToken()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedToken, token)
}

func Test_HttpTokenRetriever_Returns_Err_On_Invalid_Json(t *testing.T) {
	// Arrange
	expectedResponseBody := []byte(`<invalid-json>`)

	response := AnHttpResponse(expectedResponseBody, 200)

	mockHttpClient := new(mocks.FormPoster)
	mockResponseChecker := new(mocks.Checker)

	mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return(expectedResponseBody, nil)

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, mockHttpClient, "notused", mockResponseChecker)

	// Act
	_, err := sut.RetrieveToken()

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Failed to parse authentication token response")
}

func Test_HttpTokenRetriever_Returns_Err_On_Empty_Token_Response(t *testing.T) {
	// Arrange
	expectedResponseBody := []byte(`{"access_token": ""}`)

	response := AnHttpResponse(expectedResponseBody, 200)

	mockHttpClient := new(mocks.FormPoster)
	mockResponseChecker := new(mocks.Checker)

	mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return(expectedResponseBody, nil)

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, mockHttpClient, "notused", mockResponseChecker)

	// Act
	_, err := sut.RetrieveToken()

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Received empty authentication token")
}
