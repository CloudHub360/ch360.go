package ch360

import "net/http"

type ClassifiersClient struct {
	baseUrl       string
	requestSender HttpDoer
}

func (client *ClassifiersClient) classifiersUrl() string {
	return "/classifiers"
}

func (client *ClassifiersClient) CreateClassifier(name string) error {
	request, err := http.NewRequest("POST", client.baseUrl + client.classifiersUrl() + "/" + name, nil)

	if err != nil {
		return err
	}

	_, err = client.requestSender.Do(request)
	return err
}
