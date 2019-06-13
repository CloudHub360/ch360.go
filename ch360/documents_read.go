package ch360

import (
	"context"
	"io"
)

//go:generate mockery -name "DocumentReader"

type DocumentReader interface {
	Read(ctx context.Context, documentId string) error
	ReadResult(ctx context.Context, documentId string, mode ReadMode) (io.ReadCloser, error)
}

func (client *DocumentsClient) Read(ctx context.Context, documentId string) error {
	response, err := newRequest(ctx, "PUT", client.baseUrl+"/documents/"+documentId+"/reads", nil).
		issue(client.requestSender)

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
	ReadWvdoc
)

func (mode ReadMode) IsBinary() bool {
	return mode != ReadText
}

var readModeHeaders = map[ReadMode]string{
	ReadPDF:   "application/pdf",
	ReadText:  "text/plain",
	ReadWvdoc: "application/vnd.waives.resultformats.read+zip",
}

func (client *DocumentsClient) ReadResult(ctx context.Context,
	documentId string, mode ReadMode) (io.ReadCloser, error) {
	headers := map[string]string{
		"Accept": readModeHeaders[mode],
	}

	response, err := newRequest(ctx, "GET", client.baseUrl+"/documents/"+documentId+"/reads", nil).
		withHeaders(headers).
		issue(client.requestSender)

	if err != nil {
		return nil, err
	}

	return response.Body, err
}
