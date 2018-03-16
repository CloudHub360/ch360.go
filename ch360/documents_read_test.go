package ch360_test

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
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

var exampleReadsHttpResponse *http.Response

type DocumentsClientReadSuite struct {
	suite.Suite
	sut          *ch360.DocumentsClient
	httpClient   *mocks.HttpDoer
	fileContents []byte
	documentId   string
	ctx          context.Context
}

func (suite *DocumentsClientReadSuite) SetupTest() {
	exampleReadsHttpResponse = anHttpResponse([]byte(exampleReadsResponseBody))

	suite.httpClient = new(mocks.HttpDoer)
	suite.sut = ch360.NewDocumentsClient(apiUrl, suite.httpClient)

	suite.fileContents = generators.Bytes()
	suite.documentId = generators.String("documentId")

	suite.ctx = context.Background()
}

func TestDocumentsClientReadSuiteRunner(t *testing.T) {
	suite.Run(t, new(DocumentsClientReadSuite))
}

func (suite *DocumentsClientReadSuite) AssertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
}

func (suite *DocumentsClientReadSuite) request() *http.Request {
	require.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	require.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func (suite *DocumentsClientReadSuite) requestHeader(name string) string {
	headerMap := suite.request().Header

	suite.Require().Len(headerMap[name], 1)

	return headerMap[name][0]
}

func (suite *DocumentsClientReadSuite) Test_Correct_Url_Is_Called() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(exampleReadsHttpResponse, nil)

	suite.sut.Read(suite.ctx, suite.documentId)

	suite.AssertRequestIssued("PUT", apiUrl+"/documents/"+suite.documentId+"/reads")
}

func (suite *DocumentsClientReadSuite) Test_Response_Body_Is_Returned_From_HttpClient() {
	suite.httpClient.
		On("Do", mock.Anything).
		Return(exampleReadsHttpResponse, nil)

	data, _ := suite.sut.ReadResult(suite.ctx, suite.documentId, ch360.ReadPDF)

	suite.Assert().Equal(exampleReadsHttpResponse.Body, data)
}

func (suite *DocumentsClientReadSuite) Test_Err_Is_Returned_From_HttpClient_When_Requesting_Result() {
	expectedErr := errors.New("request error")
	suite.httpClient.
		On("Do", mock.Anything).
		Return(nil, expectedErr)

	_, receivedErr := suite.sut.ReadResult(suite.ctx, suite.documentId, ch360.ReadPDF)

	suite.Assert().Equal(expectedErr, receivedErr)
}

func (suite *DocumentsClientReadSuite) Test_Err_Is_Returned_From_HttpClient_When_Performing_Read() {
	expectedErr := errors.New("request error")
	suite.httpClient.
		On("Do", mock.Anything).
		Return(nil, expectedErr)

	receivedErr := suite.sut.Read(suite.ctx, suite.documentId)

	suite.Assert().Equal(expectedErr, receivedErr)
}

func (suite *DocumentsClientReadSuite) Test_Correct_Accept_Header_Is_Set_When_Requesting_Result() {
	var testdata = []struct {
		mode           ch360.ReadMode
		expectedHeader string
	}{
		{
			mode:           ch360.ReadPDF,
			expectedHeader: "application/pdf",
		},
		{
			mode:           ch360.ReadText,
			expectedHeader: "text/plain",
		},
	}

	for _, td := range testdata {
		suite.SetupTest()
		suite.httpClient.
			On("Do", mock.Anything).
			Return(AnHttpResponse(generators.Bytes()), nil)

		suite.sut.ReadResult(suite.ctx, suite.documentId, td.mode)

		sentAcceptHeader := suite.requestHeader("Accept")
		suite.Assert().Equal(td.expectedHeader, sentAcceptHeader)
	}
}

var exampleReadsResponseBody = `{
	"_links": {
		"self": {
			"href": "/documents/fdfwscb3SUKKqLc_5Cazhg/reads",
			"method": "PUT"
		},
		"parent": {
			"href": "/documents/fdfwscb3SUKKqLc_5Cazhg"
		}
	}
}`
