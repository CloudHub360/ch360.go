package ch360_test

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/net/mocks"
	"github.com/waives/surf/test/generators"
	"net/http"
	"testing"
)

type DocumentsClientExtractSuite struct {
	suite.Suite
	sut           *ch360.DocumentsClient
	httpClient    *mocks.HttpDoer
	fileContents  *bytes.Buffer
	documentId    string
	extractorName string
}

func (suite *DocumentsClientExtractSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.fileContents = bytes.NewBuffer(generators.Bytes())

	suite.sut = ch360.NewDocumentsClient(apiUrl, suite.httpClient)

	suite.documentId = generators.String("documentId")
	suite.extractorName = generators.String("extractorName")
}

func TestExtractDocumentSuiteRunner(t *testing.T) {
	suite.Run(t, new(DocumentsClientExtractSuite))
}

func (suite *DocumentsClientExtractSuite) Test_Extract_Returns_Results() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse([]byte(exampleExtractDocumentResponse)), nil)

	extractionResult, err := suite.sut.Extract(context.Background(), suite.documentId, suite.extractorName)

	assert.Nil(suite.T(), err)
	require.Equal(suite.T(), 1, len(extractionResult.FieldResults))
	fieldResult := extractionResult.FieldResults[0]
	assert.Equal(suite.T(), "Amount", fieldResult.FieldName)
	assert.Equal(suite.T(), "$5.50", fieldResult.Result.Text)
}

func (suite *DocumentsClientExtractSuite) Test_Extract_Deserialises_Null_Result_Property_To_Nil_Field_Result() {
	// This is a regression test for issue #85
	// If a field in the extraction response is empty and has a null result (i.e. "result": null)
	// then when it is deserialised to a FieldResult, check its Result property is nil rather than
	// an empty struct.
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse([]byte(exampleExtractDocumentResponseWithNullResult)), nil)

	extractionResult, err := suite.sut.Extract(context.Background(), suite.documentId, suite.extractorName)

	assert.Nil(suite.T(), err)
	require.Equal(suite.T(), 1, len(extractionResult.FieldResults))
	fieldResult := extractionResult.FieldResults[0]
	assert.Nil(suite.T(), fieldResult.Result)
}

func (suite *DocumentsClientExtractSuite) Test_Extract_Returns_Error_If_Response_Cannot_Be_Parsed() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse([]byte("<invalid-json>")), nil)

	_, err := suite.sut.Extract(context.Background(), suite.documentId, suite.extractorName)

	assert.NotNil(suite.T(), err)
}

func (suite *DocumentsClientExtractSuite) Test_Extract_Returns_Error_From_Http_Client() {
	expectedErr := errors.New("simulated error")
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse([]byte("")), expectedErr)

	_, receivedErr := suite.sut.Extract(context.Background(), suite.documentId, suite.extractorName)

	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *DocumentsClientExtractSuite) Test_ExtractForRedaction_Returns_Error_From_Http_Client() {
	expectedErr := errors.New("simulated error")
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse([]byte("")), expectedErr)

	_, receivedErr := suite.sut.ExtractForRedaction(context.Background(), suite.documentId, suite.extractorName)

	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *DocumentsClientExtractSuite) Test_ExtractForRedaction_Issues_Correct_Request() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse(generators.Bytes()), nil)

	_, _ = suite.sut.ExtractForRedaction(context.Background(), suite.documentId, suite.extractorName)

	expectedUrl := apiUrl + "/documents/" + suite.documentId + "/extract/" + suite.extractorName
	actualAcceptHeader := suite.request().Header.Get("Accept")
	expectedRequestHeader := "application/vnd.waives.requestformats.redact+json"
	suite.assertRequestIssued("POST", expectedUrl)
	assert.Equal(suite.T(), expectedRequestHeader, actualAcceptHeader)
}

func (suite *DocumentsClientExtractSuite) Test_ExtractForRedaction_Returns_Correct_Result_From_Response() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse([]byte(exampleExtractDocumentForRedactionResult)), nil)

	actualResult, actualErr := suite.sut.ExtractForRedaction(context.Background(),
		suite.documentId, suite.extractorName)

	assert.NoError(suite.T(), actualErr)
	assert.Equal(suite.T(), 2, len(actualResult.Marks))
	assert.Equal(suite.T(), 1, len(actualResult.Bookmarks))
}

func (suite *DocumentsClientExtractSuite) assertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
}

func (suite *DocumentsClientExtractSuite) request() *http.Request {
	require.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	require.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

var exampleExtractDocumentResponseWithNullResult = `{
	"field_results": [
		{
			"field_name": "Amount",
			"rejected": false,
			"reject_reason": "Empty",
			"result": null,
			"alternative_results": null,
			"tabular_results": null
		}
	],
	"page_sizes": {
		"page_count": 1,
		"pages": [
			{
				"page_number": 1,
				"width": 611.0,
				"height": 1008.0
			}
		]
	}
}
`

var exampleExtractDocumentResponse = `{
	"field_results": [
		{
			"field_name": "Amount",
			"rejected": false,
			"reject_reason": "None",
			"result": {
				"text": "$5.50",
				"value": null,
				"rejected": false,
				"reject_reason": "None",
				"proximity_score": 100.0,
				"match_score": 100.0,
				"text_score": 100.0,
				"areas": [
					{
						"top": 558.7115,
						"left": 276.48,
						"bottom": 571.1989,
						"right": 298.58,
						"page_number": 1
					}
				]
			},
			"alternative_results": null,
			"tabular_results": null
		}
	],
	"page_sizes": {
		"page_count": 1,
		"pages": [
			{
				"page_number": 1,
				"width": 611.0,
				"height": 1008.0
			}
		]
	}
}
`

var exampleExtractDocumentForRedactionResult = `{
  "marks": [
    {
      "name": "Name",
      "area": {
        "top": 143.447037,
        "left": 90.0,
        "bottom": 156.081787,
        "right": 151.699951,
        "page_number": 1
      }
    },
    {
      "name": "Name",
      "area": {
        "top": 471.09845,
        "left": 53.76,
        "bottom": 483.415741,
        "right": 110.959984,
        "page_number": 1
      }
    }
  ],
  "apply_marks": true,
  "bookmarks": [
    {
      "text": "Name",
      "page_number": 1
    }
  ]
}`
