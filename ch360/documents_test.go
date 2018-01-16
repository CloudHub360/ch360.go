package ch360_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
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
	sut            *ch360.DocumentsClient
	httpClient     *mocks.HttpDoer
	fileContents   []byte
	documentId     string
	classifierName string
}

var exampleCreateDocHttpResponse *http.Response
var exampleClassifyDocHttpResponse *http.Response
var exampleGetAllDocsHttpResponse *http.Response

func (suite *DocumentsClientSuite) SetupTest() {
	exampleCreateDocHttpResponse = AnHttpResponse([]byte(exampleCreateDocumentResponse))
	exampleClassifyDocHttpResponse = AnHttpResponse([]byte(exampleClassifyDocumentResponse))
	exampleGetAllDocsHttpResponse = AnHttpResponse([]byte(exampleGetAllDocumentsResponse))

	suite.httpClient = new(mocks.HttpDoer)

	suite.sut = ch360.NewDocumentsClient(apiUrl, suite.httpClient)

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
	suite.httpClient.On("Do", mock.Anything).Return(exampleCreateDocHttpResponse, nil)

	suite.sut.Create(context.Background(), suite.fileContents)

	suite.AssertRequestIssued("POST", apiUrl+"/documents")
	suite.AssertRequestHasBody(suite.fileContents)
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Returns_DocumentId() {
	suite.httpClient.On("Do", mock.Anything).Return(exampleCreateDocHttpResponse, nil)

	documentId, err := suite.sut.Create(context.Background(), suite.fileContents)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "exampleDocumentId", documentId)
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Returns_Error_From_Sender() {
	expectedErr := errors.New("simulated error")

	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	documentId, err := suite.sut.Create(context.Background(), suite.fileContents)
	assert.Equal(suite.T(), "", documentId)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_CreateDocument_Returns_Error_If_DocumentId_Cannot_Be_Parsed_From_Response() {
	expectedErr := errors.New("Could not retrieve document ID from Create Document response")
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte("")), nil)

	documentId, err := suite.sut.Create(context.Background(), suite.fileContents)
	assert.Equal(suite.T(), "", documentId)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_DeleteDocument_Issues_Delete_Document_Request() {
	suite.httpClient.On("Do", mock.Anything).Return(nil, nil)

	suite.sut.Delete(context.Background(), suite.documentId)

	suite.AssertRequestIssued("DELETE", apiUrl+"/documents/"+suite.documentId)
}

func (suite *DocumentsClientSuite) Test_DeleteDocument_Returns_Error_From_Sender() {
	expectedErr := errors.New("simulated error")

	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	err := suite.sut.Delete(context.Background(), suite.documentId)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Issues_Classify_Document_Request() {
	suite.httpClient.On("Do", mock.Anything).Return(exampleClassifyDocHttpResponse, nil)

	suite.sut.Classify(context.Background(), suite.documentId, suite.classifierName)

	suite.AssertRequestIssued("POST", apiUrl+"/documents/"+suite.documentId+"/classify/"+suite.classifierName)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_Error_From_Sender() {
	expectedErr := errors.New("simulated error")
	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	classificationResult, err := suite.sut.Classify(context.Background(), suite.documentId, suite.classifierName)
	assert.Nil(suite.T(), classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Returns_Document_Type() {
	suite.httpClient.On("Do", mock.Anything).Return(exampleClassifyDocHttpResponse, nil)

	classificationResult, err := suite.sut.Classify(context.Background(), suite.documentId, suite.classifierName)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Assignment of Deed of Trust", classificationResult.DocumentType)
}

func (suite *DocumentsClientSuite) Test_ClassifyDocument_Indicates_Confidence_Of_Result() {
	suite.httpClient.On("Do", mock.Anything).Return(exampleClassifyDocHttpResponse, nil)

	classificationResult, err := suite.sut.Classify(context.Background(), suite.documentId, suite.classifierName)

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
	expectedErr := errors.New("Could not retrieve document type from Classify response")
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte("")), nil)

	classificationResult, err := suite.sut.Classify(context.Background(), suite.documentId, suite.classifierName)

	assert.Nil(suite.T(), classificationResult)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *DocumentsClientSuite) Test_GetAll_Issues_Get_All_Documents_Request() {
	suite.httpClient.On("Do", mock.Anything).Return(exampleGetAllDocsHttpResponse, nil)

	suite.sut.GetAll(context.Background())

	suite.AssertRequestIssued("GET", apiUrl+"/documents")
}

func (suite *DocumentsClientSuite) Test_GetAll_Documents_Returns_List_Of_Documents() {
	suite.httpClient.On("Do", mock.Anything).Return(exampleGetAllDocsHttpResponse, nil)

	docs, _ := suite.sut.GetAll(context.Background())

	assert.Equal(suite.T(), 2, len(docs))
	assert.Equal(suite.T(), "yOq34IxGWk-_kAfQUdlcbw", docs[0].Id)
}

func (suite *DocumentsClientSuite) Test_GetAll_Documents_Returns_Error_From_Http_Client() {
	expectedErr := errors.New("expected")
	suite.httpClient.On("Do", mock.Anything).Return(nil, expectedErr)

	_, receivedErr := suite.sut.GetAll(context.Background())

	assert.Equal(suite.T(), expectedErr, receivedErr)
}

func (suite *DocumentsClientSuite) Test_GetAll_Documents_Returns_Error_If_Response_Cannot_Be_Parsed() {
	expectedErr := errors.New("Could not parse response")
	suite.httpClient.On("Do", mock.Anything).Return(AnHttpResponse([]byte("<invalid-json>")), nil)

	_, receivedErr := suite.sut.GetAll(context.Background())

	assert.Equal(suite.T(), expectedErr, receivedErr)
}

var exampleCreateDocumentResponse = `{
	"id": "exampleDocumentId"
}`

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

var exampleGetAllDocumentsResponse = `{
	"documents": [
		{
			"id": "yOq34IxGWk-_kAfQUdlcbw",
			"_links": {
				"document:classify": {
					"href": "/documents/yOq34IxGWk-_kAfQUdlcbw/classify/{classifier_name}",
					"templated": true
				},
				"self": {
					"href": "/documents/yOq34IxGWk-_kAfQUdlcbw"
				}
			},
			"_embedded": {
				"files": [
					{
						"id": "mDIr2ixUs0O5n3wEcMOB5g",
						"file_type": "PDF:PDFMisc",
						"size": 112449,
						"sha256": "dffe7ff587dfbd7c1dca771529c802994d5dad432986e1aaeae189b9acd40753"
					}
				]
			}
		},
		{
			"id": "7rRf0hWbHUaGua7oDszMpQ",
			"_links": {
				"document:classify": {
					"href": "/documents/7rRf0hWbHUaGua7oDszMpQ/classify/{classifier_name}",
					"templated": true
				},
				"self": {
					"href": "/documents/7rRf0hWbHUaGua7oDszMpQ"
				}
			},
			"_embedded": {
				"files": [
					{
						"id": "_QKD9ONwoU-5KueLxvcKrQ",
						"file_type": "PDF:PDFMisc",
						"size": 112449,
						"sha256": "dffe7ff587dfbd7c1dca771529c802994d5dad432986e1aaeae189b9acd40753"
					}
				]
			}
		}
	]
}`
