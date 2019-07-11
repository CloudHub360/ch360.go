package ch360_test

import (
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ResponseCheckingDoerSuite struct {
	suite.Suite
	sut        *ch360.ResponseCheckingDoer
	underlying *mocks.HttpDoer
	checker    *mocks.ResponseChecker
}

func (suite *ResponseCheckingDoerSuite) SetupTest() {
	suite.underlying = &mocks.HttpDoer{}
	suite.checker = &mocks.ResponseChecker{}
	suite.sut = ch360.NewResponseCheckingdoer(suite.checker, suite.underlying)
}

func TestResponseCheckingDoerSuiteRunner(t *testing.T) {
	suite.Run(t, new(ResponseCheckingDoerSuite))
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Calls_Underlying() {
	// Arrange
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.checker.On("CheckForErrors", mock.Anything).Return(nil)

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
	suite.checker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Calls_Checker() {
	// Arrange
	response := http.Response{}
	suite.underlying.On("Do", mock.Anything).Return(&response, nil)
	suite.checker.On("CheckForErrors", mock.Anything).Return(nil)

	// Act
	suite.sut.Do(&http.Request{})

	// Assert
	suite.checker.AssertNumberOfCalls(suite.T(), "CheckForErrors", 1)
	suite.checker.AssertCalled(suite.T(), "CheckForErrors", &response)
}

func (suite *ResponseCheckingDoerSuite) Test_ResponseCheckingDoer_Returns_Err_From_Checker() {
	// Arrange
	expectedErr := errors.New("an error")
	suite.underlying.On("Do", mock.Anything).Return(nil, nil)
	suite.checker.On("CheckForErrors", mock.Anything).Return(expectedErr)

	// Act
	_, receivedErr := suite.sut.Do(&http.Request{})

	// Assert
	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func Test_DetailedErrorResponse_Error_Returns_Detail_Message(t *testing.T) {
	fixtures := []struct {
		title       string
		detail      string
		expectedErr string
	}{
		{
			title:       "",
			detail:      "detail message",
			expectedErr: "detail message",
		}, {
			title:       "title",
			detail:      "detail message",
			expectedErr: "detail message",
		}, {
			title:       "title",
			expectedErr: "title",
		},
	}

	for _, fixture := range fixtures {

		sut := net.DetailedErrorResponse{
			Title:  fixture.title,
			Detail: fixture.detail,
		}

		assert.Equal(t, fixture.expectedErr, sut.Error())
	}
}
