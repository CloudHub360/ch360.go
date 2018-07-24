package ch360

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/net"
	"io"
	"net/http"
)

type ClassifiersClient struct {
	baseUrl       string
	requestSender net.HttpDoer
}

func NewClassifiersClient(baseUrl string, requestSender net.HttpDoer) *ClassifiersClient {
	return &ClassifiersClient{
		baseUrl:       baseUrl,
		requestSender: requestSender,
	}
}

type Classifier struct {
	Name string
}

type ClassifierList []Classifier

func (client *ClassifiersClient) issueRequest(ctx context.Context, method string, classifierName string) (*http.Response, error) {
	return client.issueRequestWith(ctx, method, classifierName, nil, nil)
}

func (client *ClassifiersClient) issueRequestWith(ctx context.Context, method string,
	suffix string,
	body io.Reader,
	headers map[string]string) (*http.Response, error) {

	request, err := http.NewRequest(method,
		client.baseUrl+"/classifiers/"+suffix,
		body)

	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	return client.requestSender.Do(request)
}

func (client *ClassifiersClient) Create(ctx context.Context, name string) error {
	_, err := client.issueRequest(ctx, "POST", name)

	return err
}

func (client *ClassifiersClient) Upload(ctx context.Context, name string, contents io.Reader) error {
	headers := map[string]string{
		"Content-Type": "application/vnd.waives.classifier+zip",
	}
	_, err := client.issueRequestWith(ctx, "POST", name, contents, headers)

	return err
}

func (client *ClassifiersClient) Delete(ctx context.Context, name string) error {
	_, err := client.issueRequest(ctx, "DELETE", name)

	return err
}

type TrainClassifierRequest struct {
	ClassifierName string
	SamplesFile    string
}

func (client *ClassifiersClient) Train(ctx context.Context, name string, samplesArchive io.Reader) error {
	headers := map[string]string{
		"Content-Type": "application/zip",
	}
	_, err := client.issueRequestWith(ctx, "POST", name+"/samples", samplesArchive, headers)

	return err
}

func (client *ClassifiersClient) GetAll(ctx context.Context) (ClassifierList, error) {

	response, err := client.issueRequest(ctx, "GET", "")

	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return nil, err
	}

	var classifiersResponse struct {
		Classifiers []Classifier
	}
	err = json.Unmarshal(buf.Bytes(), &classifiersResponse)

	if err != nil {
		return nil, err
	}

	return classifiersResponse.Classifiers, nil
}

func (classifiers ClassifierList) Contains(item string) bool {
	for _, b := range classifiers {
		if b.Name == item {
			return true
		}
	}
	return false
}

func (classifiers ClassifierList) Any() bool {
	return len(classifiers) > 0
}
