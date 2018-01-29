package ch360_test

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

type ExtractorsClientSuite struct {
	suite.Suite
	sut             *ch360.ExtractorsClient
	httpClient      *mocks.HttpDoer
	extractorName   string
	extractorConfig *bytes.Buffer
}

func (suite *ExtractorsClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut = ch360.NewExtractorsClient(apiUrl, suite.httpClient)
	suite.extractorName = "extractor-name"
	suite.extractorConfig = &bytes.Buffer{}
}

func TestExtractorsClientSuiteRunner(t *testing.T) {
	suite.Run(t, new(ExtractorsClientSuite))
}

func (suite *ExtractorsClientSuite) request() *http.Request {
	assert.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	assert.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func (suite *ExtractorsClientSuite) AssertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
}

func (suite *ExtractorsClientSuite) ClearExpectedCalls() {
	suite.httpClient.ExpectedCalls = nil
}

func (suite *ExtractorsClientSuite) Test_CreateExtractor_Issues_Create_Extractor_Request() {
	// Arrange
	suite.extractorConfig.Write([]byte("some bytes"))
	suite.extractorConfig.Reset()

	// Act
	suite.sut.Create(suite.extractorName, suite.extractorConfig)

	// Assert
	suite.AssertRequestIssued("POST", apiUrl+"/extractors/"+suite.extractorName)
}

func (suite *ExtractorsClientSuite) Test_DeleteExtractor_Issues_Delete_Extractor_Request() {
	// Act
	suite.sut.Delete(suite.extractorName)

	// Assert
	suite.AssertRequestIssued("DELETE", apiUrl+"/extractors/"+suite.extractorName)
}

func (suite *ExtractorsClientSuite) Test_GetAll_Issues_Get_All_Extractors_Request() {
	// Arrange
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(
		AnHttpResponse([]byte("{}")),
		nil)

	// Act
	suite.sut.GetAll()

	// Assert
	suite.AssertRequestIssued("GET", apiUrl+"/extractors/")
}

func (suite *ExtractorsClientSuite) Test_GetAll_Returns_List_Of_Extractors() {
	// Arrange
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(
		AnHttpResponse([]byte(exampleGetExtractorsResponse)),
		nil)

	// Act
	extractors, _ := suite.sut.GetAll()

	// Assert
	assert.Equal(suite.T(), AListOfExtractors("my-extractor", "amount"), extractors)
}

func AListOfExtractors(names ...string) ch360.ExtractorList {
	var expected ch360.ExtractorList

	for _, name := range names {
		expected = append(expected, ch360.Extractor{name})
	}

	return expected
}

func anHttpResponse(body []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
}

var exampleGetExtractorsResponse = `{
	"extractors": [
		{
			"name": "my-extractor",
			"_links": {
				"self": {
					"href": "/extractors/my-extractor"
				},
				"extractor:get": {
					"href": "/extractors/my-extractor"
				}
			}
		},
		{
			"name": "amount",
			"_links": {
				"self": {
					"href": "/extractors/amount"
				},
				"extractor:get": {
					"href": "/extractors/amount"
				}
			}
		}
	]
}`
