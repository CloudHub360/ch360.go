package auth

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"bytes"
	"github.com/CloudHub360/ch360.go/response"
	"errors"
	"github.com/stretchr/testify/suite"
	"net/url"
)

var fakeClientId = "fake-client-id"
var fakeClientSecret = "fake-client-secret"

func AnHttpResponse(body []byte, status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
	}
}

type HttpTokenRetrieverSuite struct {
	suite.Suite
	sut                 *HttpTokenRetriever
	mockHttpClient      *mocks.FormPoster
	mockResponseChecker *mocks.Checker
	validTokenValue		string
	validTokenBody      string
	validTokenResponse	*http.Response
}

func (suite *HttpTokenRetrieverSuite) SetupTest() {
	suite.mockHttpClient = new(mocks.FormPoster)
	suite.mockResponseChecker = new(mocks.Checker)
	suite.sut = NewHttpTokenRetriever(fakeClientId, fakeClientSecret, suite.mockHttpClient, "notused", suite.mockResponseChecker)
	suite.validTokenValue = `tokenvalue`
	suite.validTokenBody = `{"access_token": "` + suite.validTokenValue + `"}`
	suite.validTokenResponse = AnHttpResponse([]byte(suite.validTokenBody), 200)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(HttpTokenRetrieverSuite))
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Sends_Client_Id_And_Secret() {
	// Arrange
	suite.mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return([]byte(suite.validTokenBody), nil)

	// Act
	suite.sut.RetrieveToken()

	// Assert
	suite.mockHttpClient.AssertCalled(suite.T(), "PostForm", mock.Anything, mock.Anything)
	receivedFormData := (suite.mockHttpClient.Calls[0].Arguments[1]).(url.Values)
	assert.Equal(suite.T(), []string{fakeClientId}, receivedFormData["client_id"])
	assert.Equal(suite.T(), []string{fakeClientSecret}, receivedFormData["client_secret"])
}

func Test_HttpTokenRetriever_Returns_Error_On_HttpClient_Error(t *testing.T) {
	// Arrange
	tokenGetter := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, &http.Client{}, "http://invalid-url:-1", &response.ErrorChecker{})

	// Act
	_, err := tokenGetter.RetrieveToken()

	// Assert
	assert.NotNil(t, err)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Passes_Response_To_Checker() {
	// Arrange
	suite.mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return([]byte(suite.validTokenBody), nil)

	// Act
	suite.sut.RetrieveToken()

	// Assert
	suite.mockResponseChecker.AssertCalled(suite.T(), "Check", suite.validTokenResponse, 200)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Error_On_ResponseChecker_Error() {
	// Arrange
	response := AnHttpResponse(nil, 200)

	suite.mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return(nil, errors.New("An error"))

	// Act
	_, err := suite.sut.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "An error")
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Parses_Token_Response() {
	// Arrange
	suite.mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return([]byte(suite.validTokenBody), nil)

	// Act
	token, err := suite.sut.RetrieveToken()

	// Assert
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.validTokenValue, token)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Err_On_Invalid_Json() {
	// Arrange
	expectedResponseBody := `<invalid-json>`
	response := AnHttpResponse([]byte(expectedResponseBody), 200)

	suite.mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return([]byte(expectedResponseBody), nil)

	// Act
	_, err := suite.sut.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "Failed to parse authentication token response")
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Err_On_Empty_Token_Response() {
	// Arrange
	expectedResponseBody := `{"access_token": ""}`
	response := AnHttpResponse([]byte(expectedResponseBody), 200)

	suite.mockHttpClient.On("PostForm", mock.Anything, mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("Check", mock.Anything, mock.Anything).Return([]byte(expectedResponseBody), nil)

	// Act
	_, err := suite.sut.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "Received empty authentication token")
}