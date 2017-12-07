package ch360

import "net/http"

type ClassifiersClient struct {
	baseUrl       string
	requestSender HttpDoer
}

func (client *ClassifiersClient) issueRequest(method string, classifierName string) error {
	request, err := http.NewRequest(method,
		client.baseUrl + "/classifiers/" + classifierName,
		nil)

	if err != nil {
		return err
	}

	_, err = client.requestSender.Do(request)

	return err
}

func (client *ClassifiersClient) CreateClassifier(name string) error {
	return client.issueRequest("POST", name)
}

func (client *ClassifiersClient) DeleteClassifier(name string) error {
	return client.issueRequest("DELETE", name)
}