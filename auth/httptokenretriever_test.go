package auth_test

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/auth"
	"github.com/waives/surf/net"
	mocknet "github.com/waives/surf/net/mocks"
	"github.com/waives/surf/test/generators"
	"io/ioutil"
	"net/http"
	"testing"
)

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
	mockResponseChecker *mocknet.ResponseChecker
	validTokenValue     string
	validTokenBody      string
	validTokenResponse  *http.Response

	clientId     string
	clientSecret string
}

func (suite *HttpTokenRetrieverSuite) SetupTest() {
	suite.mockHttpClient = new(mocknet.HttpDoer)
	suite.mockResponseChecker = new(mocknet.ResponseChecker)
	suite.sut = auth.NewHttpTokenRetriever(suite.mockHttpClient, "notused", suite.mockResponseChecker)

	suite.validTokenValue = `tokenvalue`
	suite.validTokenBody = `{
		"access_token": "` + suite.validTokenValue + `",
		"expires_in": 86400
	}`
	suite.validTokenResponse = AnHttpResponse([]byte(suite.validTokenBody))

	suite.clientId = generators.String("client-id")
	suite.clientSecret = generators.String("client-secret")
}

func TestSuiteRunner(t *testing.T) {
	suite.Run(t, new(HttpTokenRetrieverSuite))
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Sends_Client_Id_And_Secret() {
	// Arrange
	suite.mockHttpClient.On("Do", mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything, mock.Anything).Return(nil)

	// Act
	suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

	// Assert
	suite.mockHttpClient.AssertCalled(suite.T(), "Do", mock.Anything)
	assert_Request_Includes_Client_Id_And_Secret(suite.T(),
		(suite.mockHttpClient.Calls[0].Arguments[0]).(*http.Request),
		suite.clientId,
		suite.clientSecret)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Error_On_HttpClient_Error() {
	// Arrange
	tokenGetter := auth.NewHttpTokenRetriever(&http.Client{}, "http://invalid-url:-1", &net.ErrorChecker{})

	// Act
	_, err := tokenGetter.RetrieveToken(suite.clientId, suite.clientSecret)

	// Assert
	assert.NotNil(suite.T(), err)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Passes_Response_To_Checker() {
	// Arrange
	suite.mockHttpClient.On("Do", mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

	// Assert
	suite.mockResponseChecker.AssertCalled(suite.T(), "CheckForErrors", suite.validTokenResponse)
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Returns_Error_On_ResponseChecker_Error() {
	// Arrange
	response := AnHttpResponse(nil)

	suite.mockHttpClient.On("Do", mock.Anything).Return(response, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(errors.New("An error"))

	// Act
	_, err := suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "An error")
}

func (suite *HttpTokenRetrieverSuite) Test_HttpTokenRetriever_Parses_Token_Response() {
	// Arrange
	suite.mockHttpClient.On("Do", mock.Anything).Return(suite.validTokenResponse, nil)
	suite.mockResponseChecker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	token, err := suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

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
	_, err := suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

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
	_, err := suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

	// Assert
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "Received empty authentication token")
}

func assert_Request_Includes_Client_Id_And_Secret(t *testing.T,
	receivedRequest *http.Request,
	expectedClientId string,
	expectedClientSecret string) {
	receivedRequest.ParseForm()
	assert.Equal(t, []string{expectedClientId}, receivedRequest.PostForm["client_id"])
	assert.Equal(t, []string{expectedClientSecret}, receivedRequest.PostForm["client_secret"])
}
