package ch360

import (
	"context"
	"io"
	"net/http"
)

//go:generate mockery -name "DocumentReader"

type DocumentReader interface {
	Read(ctx context.Context, documentId string) error
	ReadResult(ctx context.Context, documentId string, mode ReadMode) (io.ReadCloser, error)
}

func (client *DocumentsClient) Read(ctx context.Context, documentId string) error {
	request, err := http.NewRequest("PUT",
		client.baseUrl+"/documents/"+documentId+"/reads",
		nil)
	request = request.WithContext(ctx)

	if err != nil {
		return err
	}

	response, err := client.requestSender.Do(request)
	if err != nil {
		return err
	}
	response.Body.Close()

	return nil
}

type ReadMode int

const (
	ReadPDF ReadMode = iota
	ReadText
)

var readModeHeaders = map[ReadMode]string{
	ReadPDF:  "application/pdf",
	ReadText: "text/plain",
}

func (client *DocumentsClient) ReadResult(ctx context.Context,
	documentId string, mode ReadMode) (io.ReadCloser, error) {

	request, err := http.NewRequest("GET",
		client.baseUrl+"/documents/"+documentId+"/reads",
		nil)
	request.Header.Add("Accept", readModeHeaders[mode])
	request = request.WithContext(ctx)

	if err != nil {
		return nil, err
	}

	response, err := client.requestSender.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, err
}
