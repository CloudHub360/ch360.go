package ch360_test

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/net/mocks"
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
	modulesTemplate *ch360.ExtractorTemplate
	ctx             context.Context
}

func (suite *ExtractorsClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut = ch360.NewExtractorsClient(apiUrl, suite.httpClient)
	suite.extractorName = "extractor-name"
	suite.extractorConfig = &bytes.Buffer{}
	suite.ctx = context.Background()
	suite.modulesTemplate = aModulesTemplate()
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

func (suite *ExtractorsClientSuite) AssertRequestIssued(method string, urlPath string) requestAssertion {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
	assert.Equal(suite.T(), suite.ctx, suite.request().Context())

	return requestAssertion{
		request: suite.request(),
	}
}

func (suite *ExtractorsClientSuite) ClearExpectedCalls() {
	suite.httpClient.ExpectedCalls = nil
}

func (suite *ExtractorsClientSuite) Test_CreateExtractor_Issues_Create_Extractor_Request() {
	// Arrange
	suite.extractorConfig.Write([]byte("some bytes"))

	// Act
	suite.sut.Create(suite.ctx, suite.extractorName, suite.extractorConfig)

	// Assert
	suite.AssertRequestIssued("POST", apiUrl+"/extractors/"+suite.extractorName)
}

func (suite *ExtractorsClientSuite) Test_CreateExtractorFromModules_Issues_Create_Extractor_Request() {
	// Act
	suite.sut.CreateFromModules(suite.ctx, suite.extractorName, *suite.modulesTemplate)

	// Assert
	suite.AssertRequestIssued("POST", apiUrl+"/extractors/"+suite.extractorName).
		WithBody(suite.T(), func(actualBody []byte) {
			assert.Equal(suite.T(), modulesTemplateJson, string(actualBody))
		}).
		WithHeaders(suite.T(), map[string][]string{
			"Content-Type": {"application/json"},
		})
}

func (suite *ExtractorsClientSuite) Test_DeleteExtractor_Issues_Delete_Extractor_Request() {
	// Act
	suite.sut.Delete(suite.ctx, suite.extractorName)

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
	suite.sut.GetAll(suite.ctx)

	// Assert
	suite.AssertRequestIssued("GET", apiUrl+"/extractors")
}

func (suite *ExtractorsClientSuite) Test_GetAll_Returns_List_Of_Extractors() {
	// Arrange
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(
		AnHttpResponse([]byte(exampleGetExtractorsResponse)),
		nil)

	// Act
	extractors, _ := suite.sut.GetAll(suite.ctx)

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

const modulesTemplateJson = `{"modules":[{"id":"waives.reference_number","arguments":{"format":"[A-Z]"},"field_aliases":[{"field":"Reference Number","alias":"Code"}]},{"id":"waives.date"}]}`

func aModulesTemplate() *ch360.ExtractorTemplate {
	template, _ := ch360.NewModulesTemplateFromJson(bytes.NewBufferString(
		modulesTemplateJson))
	return template
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
