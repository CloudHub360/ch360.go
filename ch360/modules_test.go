package ch360_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ModulesClientSuite struct {
	suite.Suite
	sut             *ch360.ModulesClient
	httpClient      *mocks.HttpDoer
	ModuleName      string
	ModuleConfig    *bytes.Buffer
	modulesTemplate *bytes.Buffer
	ctx             context.Context
}

func (suite *ModulesClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut = ch360.NewModulesClient(apiUrl, suite.httpClient)
	suite.ModuleName = "Module-name"
	suite.ModuleConfig = &bytes.Buffer{}
	suite.ctx = context.Background()
}

func TestModulesClientSuiteRunner(t *testing.T) {
	suite.Run(t, new(ModulesClientSuite))
}

func (suite *ModulesClientSuite) request() *http.Request {
	assert.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	assert.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func (suite *ModulesClientSuite) AssertRequestIssued(method string, urlPath string) requestAssertion {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
	assert.Equal(suite.T(), suite.ctx, suite.request().Context())

	return requestAssertion{
		request: suite.request(),
	}
}

func (suite *ModulesClientSuite) ClearExpectedCalls() {
	suite.httpClient.ExpectedCalls = nil
}

func (suite *ModulesClientSuite) Test_GetAll_Returns_List_Of_Modules() {
	// Arrange
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(
		AnHttpResponse([]byte(exampleGetModulesResponse)),
		nil)

	// Act
	modules, _ := suite.sut.GetAll(suite.ctx)

	// Assert
	assert.Equal(suite.T(), suite.aListOfModules(), modules)
	assert.Equal(suite.T(), 2, len(modules))
}

func (suite *ModulesClientSuite) aListOfModules() ch360.ModuleList {
	var modulesResponse struct {
		Modules ch360.ModuleList
	}

	err := json.Unmarshal([]byte(exampleGetModulesResponse), &modulesResponse)

	assert.NoError(suite.T(), err)
	return modulesResponse.Modules
}

var exampleGetModulesResponse = `{
  "modules": [
    {
      "id": "waives.reference_number",
      "name": "Reference Number",
      "summary": "Identifies reference numbers in a document, matching a specified format.",
      "fields": [
        {
          "name": "Reference Number",
          "description": null
        }
      ],
      "parameters": [
        {
          "id": "keywords",
          "name": "Keywords",
          "type": "List",
          "description": "Restrict results to numbers found near these words, e.g. \"Order Number\", \"Order No.\". Multiple words/phrases are specified by comma-separating them.",
          "required": true
        },
        {
          "id": "format",
          "name": "Format",
          "type": "Regex",
          "description": "The format of Reference Numbers to find, either as a simple format or a regular expression.",
          "required": true
        }
      ]
    },
    {
      "id": "waives.currency",
      "name": "Currency",
      "summary": "Identifies the currency of the document based on the currency symbols in the document and the locale in which the document was produced.",
      "fields": [
        {
          "name": "Currency",
          "description": "The three-letter ISO 4217 currency code for the document. For the en-US locale, this is one of USD, GBP, CAD, MXN, or EUR; for the en-GB locale, this is one of GBP, EUR, or USD."
        }
      ],
      "parameters": [
        {
          "id": "locale",
          "name": "Currency Locale",
          "type": "List",
          "description": "The locale(s) in which the document was produced, determining the list of probable currencies in use in the document. For example, the Mexican Peso is a probable currency for en-US invoices, but not for en-GB invoices. Multiple values can be specified by comma-separating them. Valid values are en-GB and en-US, or a combination of the two.",
          "required": true
        }
      ]
    }
  ]
}`
