package ch360

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go/build"
	"io/ioutil"
	"net/http"
	"testing"
)

type ClassifiersClientSuite struct {
	suite.Suite
	sut            *ClassifiersClient
	httpClient     *mocks.HttpDoer
	classifierName string
}

const apiUrl = "baseUrl"

func (suite *ClassifiersClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut = &ClassifiersClient{
		requestSender: suite.httpClient,
		baseUrl:       apiUrl,
	}
	suite.classifierName = "classifier-name"
}

func TestClassifiersClientSuiteRunner(t *testing.T) {
	suite.Run(t, new(ClassifiersClientSuite))
}

func (suite *ClassifiersClientSuite) request() *http.Request {
	assert.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	assert.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func (suite *ClassifiersClientSuite) AssertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
}

func (suite *ClassifiersClientSuite) ClearExpectedCalls() {
	suite.httpClient.ExpectedCalls = nil
}

func (suite *ClassifiersClientSuite) Test_CreateClassifier_Issues_Create_Classifier_Request() {
	// Act
	suite.sut.Create(suite.classifierName)

	// Assert
	suite.AssertRequestIssued("POST", apiUrl+"/classifiers/"+suite.classifierName)

}

func (suite *ClassifiersClientSuite) Test_DeleteClassifier_Issues_Delete_Classifier_Request() {
	// Act
	suite.sut.Delete(suite.classifierName)

	// Assert
	suite.AssertRequestIssued("DELETE", apiUrl+"/classifiers/"+suite.classifierName)
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
	suite.AssertRequestIssued("GET", apiUrl+"/classifiers/")
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
	assert.Equal(suite.T(), AListOfClassifiers("classifier1", "classifier2"), classifiers)
}

func (suite *ClassifiersClientSuite) Test_Train_Issues_Add_Samples_Request() {
	// Act
	err := suite.sut.Train(
		suite.classifierName,
		build.Default.GOPATH+"/src/github.com/CloudHub360/ch360.go/test/samples.zip")

	// Assert
	assert.Nil(suite.T(), err)
	suite.AssertRequestIssued("POST", apiUrl+"/classifiers/"+suite.classifierName+"/samples")
}

func AListOfClassifiers(names ...string) ClassifierList {
	expected := make(ClassifierList, len(names))

	for index, name := range names {
		expected[index] = Classifier{name}
	}

	return expected
}

func AnHttpResponse(body []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
}

var exampleGetClassifiersResponse = `
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
