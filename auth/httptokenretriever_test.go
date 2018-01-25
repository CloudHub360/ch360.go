package auth_test

import (
	"bytes"
	"errors"
	"github.com/CloudHub360/ch360.go/auth"
	mocknet "github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/CloudHub360/ch360.go/response"
	mockresponse "github.com/CloudHub360/ch360.go/response/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

var fakeClientId = "fake-client-id"
var fakeClientSecret = "fake-client-secret"

func AnHttpResponse(body []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
	}
}

type HttpTokenRetrieverSuite struct {
	suite.Suite
	sut                 *auth.HttpTokenRetriever
	mockHttpClient      *mocknet.HttpDoer
	mockResponseChecker *mockresponse.Checker
	validTokenValue     string
	validTokenBody      string
	validTokenResponse  *http.Response
}

func (suite *HttpTokenRetrieverSuite) SetupTest() {
	suite.mockHttpClient = new(mocknet.HttpDoer)
	suite.mockResponseChecker = new(mockresponse.Checker)
	suite.sut = auth.NewHttpTokenRetriever(fakeClientId, fakeClientSecret, suite.mockHttpClient, "notused", suite.mockResponseChecker)

	suite.validTokenValue = `tokenvalue`
	suite.validTokenBody = `{
		"access_token": "` + suite.validTokenValue + `",
		"expires_in": 86400
	}`
	suite.validTokenResponse = AnHttpResponse([]byte(suite.validTokenBody))
}

func TestSuiteRunner(t *testing.T) {
	suite.Run(t, new(HttpTokenRetrieverSuite))
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Sends_Client_Id_And_Secret() {
	// Arrange
	suite.mockHttpClient.On("Do", mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything, mock.Anything).Return(nil)

	// Act
	suite.sut.RetrieveToken()

	// Assert
	suite.mockHttpClient.AssertCalled(suite.T(), "Do", mock.Anything)
	assert_Request_Includes_Client_Id_And_Secret(suite.T(), (suite.mockHttpClient.Calls[0].Arguments[0]).(*http.Request))
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Error_On_HttpClient_Error() {
	// Arrange
	tokenGetter := auth.NewHttpTokenRetriever(fakeClientId, fakeClientSecret, &http.Client{}, "http://invalid-url:-1", &response.ErrorChecker{})

	// Act
	_, err := tokenGetter.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Passes_Response_To_Checker() {
	// Arrange
	suite.mockHttpClient.On("Do", mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	suite.sut.RetrieveToken()

	// Assert
	suite.mockResponseChecker.AssertCalled(suite.T(), "CheckForErrors", suite.validTokenResponse)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Error_On_ResponseChecker_Error() {
	// Arrange
	response := AnHttpResponse(nil)

	suite.mockHttpClient.On("Do", mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(errors.New("An error"))

	// Act
	_, err := suite.sut.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "An error")
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Parses_Token_Response() {
	// Arrange
	suite.mockHttpClient.On("Do", mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	token, err := suite.sut.RetrieveToken()

	// Assert
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.validTokenValue, token.TokenString)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Err_On_Invalid_Json() {
	// Arrange
	expectedResponseBody := `<invalid-json>`
	response := AnHttpResponse([]byte(expectedResponseBody))

	suite.mockHttpClient.On("Do", mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	_, err := suite.sut.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "Failed to parse authentication token response")
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Err_On_Empty_Token_Response() {
	// Arrange
	expectedResponseBody := `{"access_token": ""}`
	response := AnHttpResponse([]byte(expectedResponseBody))

	suite.mockHttpClient.On("Do", mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	_, err := suite.sut.RetrieveToken()

	// Assert
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "Received empty authentication token")
}

func assert_Request_Includes_Client_Id_And_Secret(t *testing.T, receivedRequest *http.Request) {
	receivedRequest.ParseForm()
	assert.Equal(t, []string{fakeClientId}, receivedRequest.PostForm["client_id"])
	assert.Equal(t, []string{fakeClientSecret}, receivedRequest.PostForm["client_secret"])
}
