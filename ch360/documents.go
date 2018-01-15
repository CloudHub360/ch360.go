package ch360

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"net/http"
)

//go:generate mockery -name "DocumentCreator|DocumentDeleter|DocumentClassifier|DocumentGetter"
type DocumentCreator interface {
	Create(ctx context.Context, fileContents []byte) (string, error)
}

type DocumentDeleter interface {
	Delete(ctx context.Context, documentId string) error
}

type DocumentClassifier interface {
	Classify(ctx context.Context, documentId string, classifierName string) (*types.ClassificationResult, error)
}

type DocumentGetter interface {
	GetAll(ctx context.Context) (*GetAllDocumentsResponse, error)
}

type GetAllDocumentsResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}

type GetDocumentResponse struct {
	Id string `json:"id"`
}

type createDocumentResponse GetDocumentResponse

type DocumentsClient struct {
	baseUrl       string
	requestSender HttpDoer
}

func NewDocumentsClient(baseUrl string, httpDoer HttpDoer) *DocumentsClient {
	return &DocumentsClient{
		baseUrl:       baseUrl,
		requestSender: httpDoer,
	}
}

type CreateDocumentRequest struct {
	FileContents []byte
}

type classifyDocumentResponse struct {
	Id      string                          `json:"_id"`
	Results classifyDocumentResultsResponse `json:"classification_results"`
}

type classifyDocumentResultsResponse struct {
	DocumentType       string  `json:"document_type"`
	IsConfident        bool    `json:"is_confident"`
	RelativeConfidence float64 `json:"relative_confidence"`
}

func (client *DocumentsClient) Create(ctx context.Context, fileContents []byte) (string, error) {
	request, err := http.NewRequest("POST",
		client.baseUrl+"/documents",
		bytes.NewBuffer(fileContents))
	request = request.WithContext(ctx)

	if err != nil {
		return "", err
	}

	response, err := client.requestSender.Do(request)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return "", err
	}

	documentResponse := createDocumentResponse{}
	err = json.Unmarshal(buf.Bytes(), &documentResponse)
	if err != nil {
		return "", errors.New("Could not retrieve document ID from Create Document response")
	}

	return documentResponse.Id, nil
}

func (client *DocumentsClient) Delete(ctx context.Context, documentId string) error {
	request, err := http.NewRequest("DELETE",
		client.baseUrl+"/documents/"+documentId,
		nil)
	request = request.WithContext(ctx)

	if err != nil {
		return err
	}

	_, err = client.requestSender.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *DocumentsClient) Classify(ctx context.Context, documentId string, classifierName string) (*types.ClassificationResult, error) {
	request, err := http.NewRequest("POST",
		client.baseUrl+"/documents/"+documentId+"/classify/"+classifierName,
		nil)
	request = request.WithContext(ctx)

	if err != nil {
		return nil, err
	}

	response, err := client.requestSender.Do(request)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return nil, err
	}

	classifyDocumentResponse := classifyDocumentResponse{}
	err = json.Unmarshal(buf.Bytes(), &classifyDocumentResponse)
	if err != nil {
		return nil, errors.New("Could not retrieve document type from Classify response")
	}

	return &types.ClassificationResult{
		DocumentType:       classifyDocumentResponse.Results.DocumentType,
		IsConfident:        classifyDocumentResponse.Results.IsConfident,
		RelativeConfidence: classifyDocumentResponse.Results.RelativeConfidence,
	}, nil
}
