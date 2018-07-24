package ch360

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/net"
	"io"
	"net/http"
)

type ExtractorsClient struct {
	baseUrl       string
	requestSender net.HttpDoer
}

func NewExtractorsClient(baseUrl string, requestSender net.HttpDoer) *ExtractorsClient {
	return &ExtractorsClient{
		baseUrl:       baseUrl,
		requestSender: requestSender,
	}
}

type Extractor struct {
	Name string
}

type ExtractorList []Extractor

func (client *ExtractorsClient) issueRequest(ctx context.Context, method string, suffix string) (*http.Response, error) {
	return client.issueRequestWithBody(ctx, method, suffix, nil)
}

func (client *ExtractorsClient) issueRequestWithBody(ctx context.Context, method, suffix string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method,
		client.baseUrl+"/extractors/"+suffix,
		body)

	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	return client.requestSender.Do(request)
}

func (client *ExtractorsClient) Create(ctx context.Context, name string, config io.Reader) error {
	_, err := client.issueRequestWithBody(ctx, "POST", name, config)

	return err
}

func (client *ExtractorsClient) Delete(ctx context.Context, name string) error {
	_, err := client.issueRequest(ctx, "DELETE", name)

	if err != nil {
		return err
	}

	return nil
}

func (client *ExtractorsClient) GetAll(ctx context.Context) (ExtractorList, error) {

	response, err := client.issueRequest(ctx, "GET", "")

	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return nil, err
	}

	var extractorsResponse struct {
		Extractors []Extractor
	}
	err = json.Unmarshal(buf.Bytes(), &extractorsResponse)

	if err != nil {
		return nil, err
	}

	return extractorsResponse.Extractors, nil
}

func (e ExtractorList) Contains(item string) bool {
	for _, b := range e {
		if b.Name == item {
			return true
		}
	}
	return false
}
