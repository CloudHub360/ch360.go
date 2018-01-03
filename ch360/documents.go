package ch360

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

//go:generate mockery -name "DocumentCreatorDeleterClassifier"
type DocumentCreatorDeleterClassifier interface {
	CreateDocument(fileContents []byte) (string, error)
	DeleteDocument(documentId string) error
	ClassifyDocument(documentId string, classifierName string) (string, error)
}

type DocumentsClient struct {
	baseUrl       string
	requestSender HttpDoer
}

type CreateDocumentRequest struct {
	FileContents []byte
}

type createDocumentResponse struct {
	Id string `json:"id"`
}

type classifyDocumentResponse struct {
	Id      string                          `json:"_id"`
	Results classifyDocumentResultsResponse `json:"classification_results"`
}

type classifyDocumentResultsResponse struct {
	DocumentType string `json:"document_type"`
}

//TODO: Return domain object with links to ClassifyDocument & DeleteDocument urls
func (client *DocumentsClient) CreateDocument(fileContents []byte) (string, error) {
	request, err := http.NewRequest("POST",
		client.baseUrl+"/documents",
		bytes.NewBuffer(fileContents))

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

func (client *DocumentsClient) DeleteDocument(documentId string) error {
	request, err := http.NewRequest("DELETE",
		client.baseUrl+"/documents/"+documentId,
		nil)

	if err != nil {
		return err
	}

	_, err = client.requestSender.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *DocumentsClient) ClassifyDocument(documentId string, classifierName string) (string, error) {
	request, err := http.NewRequest("POST",
		client.baseUrl+"/documents/"+documentId+"/classify/"+classifierName,
		nil)

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

	classifyDocumentResponse := classifyDocumentResponse{}
	err = json.Unmarshal(buf.Bytes(), &classifyDocumentResponse)
	if err != nil {
		return "", errors.New("Could not retrieve document type from ClassifyDocument response")
	}

	return classifyDocumentResponse.Results.DocumentType, nil
}
