package ch360

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type DocumentsClientSuite struct {
	suite.Suite
	sut            *DocumentsClient
	httpClient     *mocks.HttpDoer
	fileContents   []byte
	documentId     string
	classifierName string
}

func (suite *DocumentsClientSuite) SetupTest() {
	suite.httpClient = new(mocks.HttpDoer)
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte(exampleCreateDocumentResponse)), nil)

	suite.sut = &DocumentsClient{
		requestSender: suite.httpClient,
		baseUrl:       apiUrl,
	}

	suite.fileContents = generators.Bytes()
	suite.documentId = generators.String("documentId")
	suite.classifierName = generators.String("classifierName")
}

func TestDocumentsClientSuiteRunner(t *testing.T) {
	suite.Run(t, new(DocumentsClientSuite))
}

func (suite *DocumentsClientSuite) request() *http.Request {
	require.Len(suite.T(), suite.httpClient.Calls, 1)

	call := suite.httpClient.Calls[0]
	require.Len(suite.T(), call.Arguments, 1)

	return (call.Arguments[0]).(*http.Request)
}

func (suite *DocumentsClientSuite) AssertRequestIssued(method string, urlPath string) {
	assert.Equal(suite.T(), method, suite.request().Method)
	assert.Equal(suite.T(), urlPath, suite.request().URL.Path)
}

func (suite *DocumentsClientSuite) AssertRequestHasBody(body []byte) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(suite.request().Body)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), body, buf.Bytes())
}

func (suite *DocumentsClientSuite) ClearExpectedCalls() {
	suite.httpClient.ExpectedCalls = nil
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Issues_Create_Document_Request_With_File_Contents() {
	suite.sut.CreateDocument(context.Background(), suite.fileContents)

	suite.AssertRequestIssued("POST", apiUrl+"/documents")
	suite.AssertRequestHasBody(suite.fileContents)
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Returns_DocumentId() {
	documentId, err := suite.sut.CreateDocument(context.Background(), suite.fileContents)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "exampleDocumentId", documentId)
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Returns_Error_From_Sender() {
	expectedErr := errors.New("simulated error")
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	documentId, err := suite.sut.CreateDocument(context.Background(), suite.fileContents)
	assert.Equal(suite.T(), "", documentId)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Returns_Error_If_DocumentId_Cannot_Be_Parsed_From_Response() {
	expectedErr := errors.New("Could not retrieve document ID from Create Document response")
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte("")), nil)

	documentId, err := suite.sut.CreateDocument(context.Background(), suite.fileContents)
	assert.Equal(suite.T(), "", documentId)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_DeleteDocument_Issues_Delete_Document_Request() {
	suite.sut.DeleteDocument(context.Background(), suite.documentId)

	suite.AssertRequestIssued("DELETE", apiUrl+"/documents/"+suite.documentId)
}

func (suite *DocumentsClientSuite) Test_DeleteDocument_Returns_Error_From_Sender() {
	expectedErr := errors.New("simulated error")
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	err := suite.sut.DeleteDocument(context.Background(), suite.documentId)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Issues_Classify_Document_Request() {
	suite.sut.ClassifyDocument(context.Background(), suite.documentId, suite.classifierName)

	suite.AssertRequestIssued("POST", apiUrl+"/documents/"+suite.documentId+"/classify/"+suite.classifierName)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_Error_From_Sender() {
	expectedErr := errors.New("simulated error")
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	classificationResult, err := suite.sut.ClassifyDocument(context.Background(), suite.documentId, suite.classifierName)
	assert.Nil(suite.T(), classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_Document_Type() {
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte(exampleClassifyDocumentResponse)), nil)

	classificationResult, err := suite.sut.ClassifyDocument(context.Background(), suite.documentId, suite.classifierName)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Assignment of Deed of Trust", classificationResult.DocumentType)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_Confident_Status() {
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte(exampleClassifyDocumentResponse)), nil)

	classificationResult, err := suite.sut.ClassifyDocument(context.Background(), suite.documentId, suite.classifierName)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), true, classificationResult.IsConfident)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_RelativeConfidence() {
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte(exampleClassifyDocumentResponse)), nil)

	classificationResult, err := suite.sut.ClassifyDocument(context.Background(), suite.documentId, suite.classifierName)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1.234567, classificationResult.RelativeConfidence)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_Error_If_DocumentType_Cannot_Be_Parsed_From_Response() {
	expectedErr := errors.New("Could not retrieve document type from ClassifyDocument response")
	suite.ClearExpectedCalls()
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte("")), nil)

	classificationResult, err := suite.sut.ClassifyDocument(context.Background(), suite.documentId, suite.classifierName)

	assert.Nil(suite.T(), classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

var exampleCreateDocumentResponse = `
{
	"id": "exampleDocumentId"
}
`

var exampleClassifyDocumentResponse = `
{
	"_id": "exampleDocumentId",
	"classification_results": {
		"document_type": "Assignment of Deed of Trust",
		"relative_confidence": 1.234567,
		"is_confident": true,
		"document_type_scores": [
			{
				"document_type": "Assignment of Deed of Trust",
				"score": 61.4187
			},
			{
				"document_type": "Notice of Default",
				"score": 32.94312
			},
			{
				"document_type": "Correspondence",
				"score": 28.2860489
			},
			{
				"document_type": "Deed of Trust",
				"score": 28.0011711
			},
			{
				"document_type": "Notice of Lien",
				"score": 27.9561481
			}
		]
	}
}
`
