package ch360

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/waives/surf/net"
	"io"
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

func (client *ClassifiersClient) Create(ctx context.Context, name string) error {
	_, err := newRequest(ctx, "POST", client.baseUrl+"/classifiers/"+name, nil).
		issue(client.requestSender)

	return err
}

func (client *ClassifiersClient) Upload(ctx context.Context, name string, contents io.Reader) error {
	headers := map[string]string{
		"Content-Type": "application/vnd.waives.classifier+zip",
	}

	_, err := newRequest(ctx, "POST", client.baseUrl+"/classifiers/"+name, contents).
		withHeaders(headers).
		issue(client.requestSender)

	return err
}

func (client *ClassifiersClient) Delete(ctx context.Context, name string) error {
	_, err := newRequest(ctx, "DELETE", client.baseUrl+"/classifiers/"+name, nil).
		issue(client.requestSender)

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

	_, err := newRequest(ctx, "POST", client.baseUrl+"/classifiers/"+name+"/samples", samplesArchive).
		withHeaders(headers).
		issue(client.requestSender)

	return err
}

func (client *ClassifiersClient) GetAll(ctx context.Context) (ClassifierList, error) {
	response, err := newRequest(ctx, "GET", client.baseUrl+"/classifiers", nil).
		issue(client.requestSender)

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
