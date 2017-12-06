package ch360

import (
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ResponseCheckingDoerSuite struct {
	suite.Suite
	sut        *responseCheckingDoer
	underlying *mocks.HttpDoer
	checker    *mocks.Checker
}

func (suite *ResponseCheckingDoerSuite) SetupTest() {
	suite.underlying = &mocks.HttpDoer{}
	suite.checker = &mocks.Checker{}
	suite.sut = &responseCheckingDoer{
		checker:       suite.checker,
		wrappedSender: suite.underlying,
	}
}

func TestResponseCheckingDoerSuiteRunner(t *testing.T) {
	suite.Run(t, new(ResponseCheckingDoerSuite))
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Calls_Underlying() {
	// Arrange
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.checker.On("Check", mock.Anything).Return(nil)

	// Act
	request := http.Request{}
	suite.sut.Do(&request)

	// Assert
	suite.underlying.AssertNumberOfCalls(suite.T(), "Do", 1)
	suite.underlying.AssertCalled(suite.T(), "Do", &request)
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Returns_Err_From_Underlying() {
	// Arrange
	expectedErr := errors.New("an error")
	suite.underlying.On("Do", mock.Anything).Return(nil, expectedErr)
	suite.checker.On("Check", mock.Anything).Return(nil)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Calls_Checker() {
	// Arrange
	response := http.Response{}
	suite.underlying.On("Do", mock.Anything).Return(&response, nil)
	suite.checker.On("Check", mock.Anything).Return(nil)

	// Act
	suite.sut.Do(&http.Request{})

	// Assert
	suite.checker.AssertNumberOfCalls(suite.T(), "Check", 1)
	suite.checker.AssertCalled(suite.T(), "Check", &response)
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Returns_Err_From_Checker() {
	// Arrange
	expectedErr := errors.New("an error")
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.checker.On("Check", mock.Anything).Return(expectedErr)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}
