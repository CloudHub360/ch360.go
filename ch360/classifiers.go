package ch360

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/net"
	"io"
	"net/http"
	"os"
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

func (client *ClassifiersClient) issueRequest(method string, classifierName string) (*http.Response, error) {
	return client.issueRequestWith(method, classifierName, nil, nil)
}

func (client *ClassifiersClient) issueRequestWith(method string,
	classifierName string,
	body io.Reader,
	headers map[string]string) (*http.Response, error) {

	request, err := http.NewRequest(method,
		client.baseUrl+"/classifiers/"+classifierName,
		body)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	return client.requestSender.Do(request)
}

func (client *ClassifiersClient) Create(name string) error {
	_, err := client.issueRequest("POST", name)

	return err
}

func (client *ClassifiersClient) Upload(name string, contents io.Reader) error {
	headers := map[string]string{
		"Content-Type": "application/vnd.waives.classifier+zip",
	}
	_, err := client.issueRequestWith("POST", name, contents, headers)

	return err
}

func (client *ClassifiersClient) Delete(name string) error {
	_, err := client.issueRequest("DELETE", name)

	if err != nil {
		return err
	}

	return nil
}

type TrainClassifierRequest struct {
	ClassifierName string
	SamplesFile    string
}

func (_req *TrainClassifierRequest) Issue(client *ClassifiersClient) error {
	zip, err := os.Open(_req.SamplesFile)
	if err != nil {
		return errors.New(fmt.Sprintf("The file '%s' could not be found.", _req.SamplesFile))
	}

	request, err := http.NewRequest("POST",
		client.baseUrl+"/classifiers/"+_req.ClassifierName+"/samples",
		zip)

	request.Header.Set("Content-Type", "application/zip")

	if err != nil {
		return err
	}

	_, err = client.requestSender.Do(request)

	if err != nil {
		return err
	}

	return nil
}

func (client *ClassifiersClient) Train(name string, samplesPath string) error {
	request := &TrainClassifierRequest{
		ClassifierName: name,
		SamplesFile:    samplesPath,
	}

	return request.Issue(client)
}

func (client *ClassifiersClient) GetAll() (ClassifierList, error) {

	response, err := client.issueRequest("GET", "")

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
