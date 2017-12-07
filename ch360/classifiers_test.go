package ch360

import (
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"github.com/stretchr/testify/suite"
)

type ClassifiersClientSuite struct {
	suite.Suite
	sut        *ClassifiersClient
	httpClient *mocks.HttpDoer
	classifierName string
}

func (suite *ClassifiersClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut = &ClassifiersClient{
		requestSender: suite.httpClient,
		baseUrl:       "baseurl",
	}
	suite.classifierName = "classifier-name"
}

func TestClassifiersClientSuiteRunner(t *testing.T) {
	suite.Run(t, new (ClassifiersClientSuite))
}

func (suite *ClassifiersClientSuite) Request() *http.Request {
	assert.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	assert.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func (suite *ClassifiersClientSuite) AssertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.Request().Method)
	assert.Equal(suite.T(), urlPath, suite.Request().URL.Path)
}

func (suite *ClassifiersClientSuite) Test_CreateClassifier_Issues_Create_Classifier_Request() {
	// Act
	suite.sut.CreateClassifier(suite.classifierName)

	// Assert
	suite.AssertRequestIssued("POST", "baseurl/classifiers/"+suite.classifierName)

}

func (suite *ClassifiersClientSuite) Test_DeleteClassifier_Issues_Delete_Classifier_Request() {
	// Act
	suite.sut.DeleteClassifier(suite.classifierName)

	// Assert
	suite.AssertRequestIssued("DELETE", "baseurl/classifiers/"+suite.classifierName)
}
