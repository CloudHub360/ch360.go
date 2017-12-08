package ch360

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ClassifiersClient struct {
	baseUrl       string
	requestSender HttpDoer
}

type Classifier struct {
	Name string
}

type ClassifierList []Classifier

func (client *ClassifiersClient) issueRequest(method string, classifierName string) (*http.Response, error) {
	request, err := http.NewRequest(method,
		client.baseUrl+"/classifiers/"+classifierName,
		nil)

	if err != nil {
		return nil, err
	}

	return client.requestSender.Do(request)
}

func (client *ClassifiersClient) Create(name string) error {
	_, err := client.issueRequest("POST", name)

	return err
}

func (client *ClassifiersClient) Delete(name string) error {
	_, err := client.issueRequest("DELETE", name)

	return err
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
