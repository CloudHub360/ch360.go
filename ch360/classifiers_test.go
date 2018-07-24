package ch360_test

import (
	"bytes"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go/build"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type ClassifiersClientSuite struct {
	suite.Suite
	sut            *ch360.ClassifiersClient
	httpClient     *mocks.HttpDoer
	classifierName string
	classifierFile io.Reader
}

const apiUrl = "baseUrl"

func (suite *ClassifiersClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut = ch360.NewClassifiersClient(apiUrl, suite.httpClient)
	suite.classifierName = "classifier-name"
	suite.classifierFile, _ = os.Open("testdata/emptyclassifier.clf")
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

func (suite *ClassifiersClientSuite) AssertRequestIssuedWith(method string, urlPath string, expectedBody io.Reader, headers map[string]string) {
	suite.AssertRequestIssued(method, urlPath)

	assert.Equal(suite.T(), expectedBody, suite.request().Body)

	for k, expectedValue := range headers {
		actualValue, ok := suite.request().Header[k]

		assert.True(suite.T(), ok)
		assert.Contains(suite.T(), actualValue, expectedValue)
	}
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

func (suite *ClassifiersClientSuite) Test_UploadClassifier_Issues_Create_Classifier_Request() {
	// Act
	suite.sut.Upload(suite.classifierName, suite.classifierFile)
	expectedHeaders := map[string]string{
		"Content-Type": "application/vnd.waives.classifier+zip",
	}

	// Assert
	suite.AssertRequestIssuedWith("POST",
		apiUrl+"/classifiers/"+suite.classifierName,
		suite.classifierFile,
		expectedHeaders)

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

func (suite *ClassifiersClientSuite) Test_Train_Returns_An_Error_If_The_Sample_Path_Is_Invalid() {
	// Act
	samplesPath := build.Default.GOPATH + "/src/github.com/CloudHub360/ch360.go/test/non-existent.zip"
	err := suite.sut.Train(
		suite.classifierName,
		samplesPath)

	// Assert
	assert.Equal(suite.T(), errors.New("The file '"+samplesPath+"' could not be found."), err)
}

func AListOfClassifiers(names ...string) ch360.ClassifierList {
	expected := make(ch360.ClassifierList, len(names))

	for index, name := range names {
		expected[index] = ch360.Classifier{name}
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
