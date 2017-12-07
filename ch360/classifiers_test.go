package ch360

import (
	"github.com/CloudHub360/ch360.go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"bytes"
)

type ClassifiersClientSuite struct {
	suite.Suite
	sut            *ClassifiersClient
	httpClient     *mocks.HttpDoer
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
	suite.Run(t, new(ClassifiersClientSuite))
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

func (suite *ClassifiersClientSuite) ClearExpectedCalls() {
	suite.httpClient.ExpectedCalls = nil
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

func (suite *ClassifiersClientSuite) Test_GetAll_Issues_Get_All_Classifiers_Request() {
	// Arrange
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(
		AnHttpResponse([]byte("{}")),
		nil)

	// Act
	suite.sut.GetAll()

	// Assert
	suite.AssertRequestIssued("GET", "baseurl/classifiers/")
}

func (suite *ClassifiersClientSuite) Test_GetAll_Returns_List_Of_Classifiers() {
	// Arrange
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(
		AnHttpResponse([]byte(exampleGetClassifiersResponse)),
		nil)

	// Act
	classifiers, _ := suite.sut.GetAll()

	// Assert
	assert.Len(suite.T(), classifiers, 2)
	assert.Equal(suite.T(), Classifier{ Name: "classifier1"}  ,classifiers[0])
	assert.Equal(suite.T(), Classifier{ Name: "classifier2"}  ,classifiers[1])
}

func AnHttpResponse(body []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
}

var exampleGetClassifiersResponse =`
{
	"classifiers": [
		{
			"name": "classifier1",
			"_links": {
				"self": {
					"href": "/classifiers/classifier1"
				},
				"classifier:add_sample": {
					"href": "/classifiers/classifier1/sample/{document_type}",
					"templated": true
				},
				"classifier:add_samples_from_zip": {
					"href": "/classifiers/classifier1/samples"
				},
				"classifier:get": {
					"href": "/classifiers/classifier1"
				}
			}
		},
		{
			"name": "classifier2",
			"_links": {
				"self": {
					"href": "/classifiers/classifier2"
				},
				"classifier:add_sample": {
					"href": "/classifiers/classifier2/sample/{document_type}",
					"templated": true
				},
				"classifier:add_samples_from_zip": {
					"href": "/classifiers/classifier2/samples"
				},
				"classifier:get": {
					"href": "/classifiers/classifier2"
				}
			}
		}
	]
}
`