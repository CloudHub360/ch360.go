package ch360

import (
	"bytes"
	"errors"
	"github.com/Jeffail/gabs"
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

	jsonParsed, err := gabs.ParseJSON(buf.Bytes())

	if documentId, ok := jsonParsed.Path("id").Data().(string); ok {
		return documentId, nil
	}
	return "", errors.New("Could not retrieve document ID from Create Document response")
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

	jsonParsed, err := gabs.ParseJSON(buf.Bytes())
	var documentType string
	var ok bool

	//TODO: Return results struct
	documentType, ok = jsonParsed.Path("classification_results.document_type").Data().(string)
	if ok {
		return documentType, nil
	}

	return "", errors.New("Could not retrieve document type from ClassifyDocument Document response")
}
