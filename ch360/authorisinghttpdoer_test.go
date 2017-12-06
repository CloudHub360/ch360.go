package ch360

import (
	"errors"
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthorisingHttpDoerSuite struct {
	suite.Suite
	sut        *authorisingDoer
	underlying *mocks.HttpDoer
	retriever  *mocks.TokenRetriever
}

func (suite *AuthorisingHttpDoerSuite) SetupTest() {
	suite.underlying = &mocks.HttpDoer{}
	suite.retriever = &mocks.TokenRetriever{}
	suite.sut = &authorisingDoer{
		tokenRetriever: suite.retriever,
		wrappedSender:  suite.underlying,
	}
}

func TestAuthorisingDoerSuiteRunner(t *testing.T) {
	suite.Run(t, new(AuthorisingHttpDoerSuite))
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Calls_TokenRetriever() {
	// Arrange
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.retriever.On("RetrieveToken").Return("token", nil)

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
	suite.retriever.On("RetrieveToken").Return("", expectedErr)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Returns_Error_From_Underlying() {
	// Arrange
	expectedErr := errors.New("underlying error")
	suite.underlying.On("Do", mock.Anything).Return(nil, expectedErr)
	suite.retriever.On("RetrieveToken", mock.Anything).Return("token", nil)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *AuthorisingHttpDoerSuite) Test_AuthorisingDoer_Calls_Underlying_With_Token() {
	// Arrange
	token := "auth_token"
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.retriever.On("RetrieveToken", mock.Anything).Return(token, nil)

	// Act
	request := http.Request{}
	suite.sut.Do(&request)

	// Assert
	assert.Equal(suite.T(), "Bearer "+token, request.Header.Get("Authorization"))
	suite.underlying.AssertNumberOfCalls(suite.T(), "Do", 1)
	suite.underlying.AssertCalled(suite.T(), "Do", &request)
}
