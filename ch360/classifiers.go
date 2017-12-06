package ch360

import "net/http"

type ClassifiersClient struct {
	baseUrl string
	sender  HttpDoer
}

func (client *ClassifiersClient) classifiersUrl() string {
	return "/classifiers"
}

func (client *ClassifiersClient) CreateClassifier(name string) error {
	request, err := http.NewRequest("POST", client.baseUrl + client.classifiersUrl() + "/" + name, nil)

	if err != nil {
		return err
	}

	_, err = client.sender.Do(request)
	return err
}
