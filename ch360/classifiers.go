package ch360

import (
	"net/http"
	"bytes"
	"encoding/json"
)

type ClassifiersClient struct {
	baseUrl       string
	requestSender HttpDoer
}

type Classifier struct {
	Name string
}

func (client *ClassifiersClient) issueRequest(method string, classifierName string) (*http.Response, error) {
	request, err := http.NewRequest(method,
		client.baseUrl+"/classifiers/"+classifierName,
		nil)

	if err != nil {
		return nil, err
	}

	return client.requestSender.Do(request)
}

func (client *ClassifiersClient) CreateClassifier(name string) error {
	_, err := client.issueRequest("POST", name)

	return err
}

func (client *ClassifiersClient) DeleteClassifier(name string) error {
	_, err := client.issueRequest("DELETE", name)

	return err
}

func (client *ClassifiersClient) GetAll() ([]Classifier, error) {

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
