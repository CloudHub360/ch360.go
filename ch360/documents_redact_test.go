package ch360_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/request"
	"github.com/CloudHub360/ch360.go/net/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type DocumentsClientRedactSuite struct {
	suite.Suite
	sut          *ch360.DocumentsClient
	httpClient   *mocks.HttpDoer
	fileContents *bytes.Buffer
	documentId   string
	redactorName string
	responseBody []byte
	pdfRequest   request.RedactedPdfRequest
}

func (suite *DocumentsClientRedactSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.fileContents = bytes.NewBuffer(generators.Bytes())

	suite.sut = ch360.NewDocumentsClient(apiUrl, suite.httpClient)

	suite.documentId = generators.String("documentId")
	suite.responseBody = generators.Bytes()
	suite.pdfRequest = aRedactPdfRequest()
}

func TestRedactDocumentSuiteRunner(t *testing.T) {
	suite.Run(t, new(DocumentsClientRedactSuite))
}

func (suite *DocumentsClientRedactSuite) Test_Redact_Returns_Body_Of_Http_Response() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse(suite.responseBody), nil)

	actualResult, err := suite.sut.Redact(context.Background(), suite.documentId, suite.pdfRequest)
	actualResultBuffer := bytes.Buffer{}
	_, _ = actualResultBuffer.ReadFrom(actualResult)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.responseBody, actualResultBuffer.Bytes())
}

func (suite *DocumentsClientRedactSuite) Test_Redact_Returns_Error_From_HttpClient() {
	expectedErr := errors.New("simulatedError")
	suite.httpClient.
		On("Do", mock.Anything).
		Return(nil, expectedErr)

	_, actualErr := suite.sut.Redact(context.Background(), suite.documentId, suite.pdfRequest)

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func (suite *DocumentsClientRedactSuite) Test_Redact_Issues_Correct_Request() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(AnHttpResponse(suite.responseBody), nil)

	_, _ = suite.sut.Redact(context.Background(), suite.documentId, suite.pdfRequest)

	suite.AssertRequestIssued("POST", apiUrl+"/documents/"+suite.documentId+"/redact")
}

func (suite *DocumentsClientRedactSuite) AssertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
}

func (suite *DocumentsClientRedactSuite) request() *http.Request {
	require.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	require.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func aRedactPdfRequest() request.RedactedPdfRequest {
	req := &request.RedactedPdfRequest{}
	_ = json.Unmarshal([]byte(exampleRedactDocumentRequest), req)
	return *req
}

var exampleRedactDocumentRequest = `{
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
