package ch360

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ClassifiersClient struct {
	baseUrl       string
	requestSender HttpDoer
}

type Classifier struct {
	Name string
}

type Request interface {
	Issue(client *ClassifiersClient) error
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

type TrainClassifierRequest struct {
	ClassifierName string
	SamplesFile    string
}

func (_req *TrainClassifierRequest) Issue(client *ClassifiersClient) error {
	fmt.Printf("Adding samples from file '%s'... ", _req.SamplesFile)

	zip, err := os.Open(_req.SamplesFile)
	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}

	request, err := http.NewRequest("POST",
		client.baseUrl+"/classifiers/"+_req.ClassifierName+"/samples",
		zip)

	request.Header.Set("Content-Type", "application/zip")

	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}

	_, err = client.requestSender.Do(request)

	if err != nil {
		fmt.Println("[FAILED]")
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
