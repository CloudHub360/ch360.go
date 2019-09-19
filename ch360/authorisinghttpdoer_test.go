package ch360_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/auth"
	mockauth "github.com/waives/surf/auth/mocks"
	"github.com/waives/surf/ch360"
	mocknet "github.com/waives/surf/net/mocks"
	"github.com/waives/surf/test/generators"
	"net/http"
	"testing"
)

type AuthorisingHttpDoerSuite struct {
	suite.Suite
	sut        *ch360.AuthorisingDoer
	underlying *mocknet.HttpDoer
	retriever  *mockauth.TokenRetriever
}

func (suite *AuthorisingHttpDoerSuite) SetupTest() {
	suite.underlying = &mocknet.HttpDoer{}
	suite.retriever = &mockauth.TokenRetriever{}
	clientId := generators.String("client-id")
	clientSecret := generators.String("client-secret")
	suite.sut = ch360.NewAuthorisingDoer(suite.retriever, suite.underlying, clientId, clientSecret)
}

func TestAuthorisingDoerSuiteRunner(t *testing.T) {
	suite.Run(t, new(AuthorisingHttpDoerSuite))
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Calls_TokenRetriever() {
	// Arrange
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.retriever.On("RetrieveToken", mock.Anything, mock.Anything).Return(&auth.AccessToken{}, nil)

	// Act
	request := http.Request{}
	suite.sut.Do(&request)

	// Assert
	suite.retriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 1)
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Returns_Error_From_TokenRetriever() {
	// Arrange
	expectedErr := errors.New("retriever error")
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.retriever.On("RetrieveToken", mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Returns_Error_From_Underlying() {
	// Arrange
	expectedErr := errors.New("underlying error")
	suite.underlying.On("Do", mock.Anything).Return(nil, expectedErr)
	suite.retriever.On("RetrieveToken", mock.Anything, mock.Anything).Return(&auth.AccessToken{}, nil)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Calls_Underlying_With_Token() {
	// Arrange
	token := &auth.AccessToken{}
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.retriever.On("RetrieveToken", mock.Anything, mock.Anything).Return(token, nil)

	// Act
	request := http.Request{}
	suite.sut.Do(&request)

	// Assert
	assert.Equal(suite.T(), "Bearer "+token.TokenString, request.Header.Get("Authorization"))
	suite.underlying.AssertNumberOfCalls(suite.T(), "Do", 1)
	suite.underlying.AssertCalled(suite.T(), "Do", &request)
}
