package ch360

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/net"
	"io"
)

var TotalDocumentSlots = 10

//go:generate mockery -name "DocumentCreator|DocumentDeleter|DocumentClassifier|DocumentGetter|DocumentExtractor"
type DocumentCreator interface {
	Create(ctx context.Context, fileContents io.Reader) (string, error)
}

type DocumentExtractor interface {
	Extract(ctx context.Context, documentId string, extractorName string) (*results.ExtractionResult, error)
}

type DocumentDeleter interface {
	Delete(ctx context.Context, documentId string) error
}

type DocumentClassifier interface {
	Classify(ctx context.Context, documentId string, classifierName string) (*results.ClassificationResult, error)
}

type DocumentGetter interface {
	GetAll(ctx context.Context) (DocumentList, error)
}

type Document struct {
	Id       string
	Size     int
	Sha256   string
	FileType string
}

type getAllDocumentsResponse struct {
	Documents []getDocumentResponse `json:"documents"`
}

type getDocumentResponse struct {
	Id       string `json:"id"`
	Embedded struct {
		Files []struct {
			ID       string `json:"id"`
			FileType string `json:"file_type"`
			Size     int    `json:"size"`
			Sha256   string `json:"sha256"`
		} `json:"files"`
	} `json:"_embedded"`
}

type createDocumentResponse struct {
	Id string `json:"id"`
}

type DocumentsClient struct {
	baseUrl       string
	requestSender net.HttpDoer
}

func NewDocumentsClient(baseUrl string, httpDoer net.HttpDoer) *DocumentsClient {
	return &DocumentsClient{
		baseUrl:       baseUrl,
		requestSender: httpDoer,
	}
}

type CreateDocumentRequest struct {
	FileContents []byte
}

type classifyDocumentResponse struct {
	Id      string `json:"_id"`
	Results struct {
		DocumentType       string  `json:"document_type"`
		IsConfident        bool    `json:"is_confident"`
		RelativeConfidence float64 `json:"relative_confidence"`
		DocumentTypeScores []struct {
			DocumentType string  `json:"document_type"`
			Score        float64 `json:"score"`
		} `json:"document_type_scores"`
	} `json:"classification_results"`
}

func (client *DocumentsClient) Create(ctx context.Context, fileContents io.Reader) (string, error) {
	response, err := newRequest(ctx, "POST", client.baseUrl+"/documents", fileContents).
		issue(client.requestSender)

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
	_, err := newRequest(ctx, "DELETE", client.baseUrl+"/documents/"+documentId, nil).
		issue(client.requestSender)

	return err
}

func (client *DocumentsClient) Classify(ctx context.Context, documentId string, classifierName string) (*results.ClassificationResult, error) {
	response, err := newRequest(ctx, "POST",
		client.baseUrl+"/documents/"+documentId+"/classify/"+classifierName, nil).
		issue(client.requestSender)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

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

	var scores []results.DocumentTypeScore
	for _, score := range classifyDocumentResponse.Results.DocumentTypeScores {
		scores = append(scores, results.DocumentTypeScore{DocumentType: score.DocumentType, Score: score.Score})
	}

	return &results.ClassificationResult{
		DocumentType:       classifyDocumentResponse.Results.DocumentType,
		IsConfident:        classifyDocumentResponse.Results.IsConfident,
		RelativeConfidence: classifyDocumentResponse.Results.RelativeConfidence,
		DocumentTypeScores: scores,
	}, nil
}

func (client *DocumentsClient) Extract(ctx context.Context, documentId string, extractorName string) (*results.ExtractionResult, error) {
	response, err := newRequest(ctx, "POST",
		client.baseUrl+"/documents/"+documentId+"/extract/"+extractorName, nil).
		issue(client.requestSender)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	extractResponse := results.ExtractionResult{}
	err = json.NewDecoder(response.Body).Decode(&extractResponse)

	if err != nil {
		return nil, err
	}

	return &extractResponse, nil
}

func (client *DocumentsClient) GetAll(ctx context.Context) (DocumentList, error) {
	response, err := newRequest(ctx, "GET",
		client.baseUrl+"/documents", nil).
		issue(client.requestSender)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return nil, err
	}

	var allDocsResponse getAllDocumentsResponse
	err = json.Unmarshal(buf.Bytes(), &allDocsResponse)
	if err != nil {
		return nil, err
	}

	var docs []Document
	for _, doc := range allDocsResponse.Documents {
		file := doc.Embedded.Files[0]
		docs = append(docs, Document{
			Id:       doc.Id,
			FileType: file.FileType,
			Sha256:   file.Sha256,
			Size:     file.Size,
		})
	}

	return docs, nil
}

type DocumentList []Document
