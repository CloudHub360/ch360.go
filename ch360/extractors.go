package ch360

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ExtractorsClient struct {
	baseUrl       string
	requestSender HttpDoer
}

func NewExtractorsClient(baseUrl string, requestSender HttpDoer) *ExtractorsClient {
	return &ExtractorsClient{
		baseUrl:       baseUrl,
		requestSender: requestSender,
	}
}

type Extractor struct {
	Name string
}

type ExtractorList []Extractor

func (client *ExtractorsClient) issueRequest(method string, extractorName string) (*http.Response, error) {
	request, err := http.NewRequest(method,
		client.baseUrl+"/extractors/"+extractorName,
		nil)

	if err != nil {
		return nil, err
	}

	return client.requestSender.Do(request)
}

func (client *ExtractorsClient) Create(name string) error {
	_, err := client.issueRequest("POST", name)

	return err
}

func (client *ExtractorsClient) Delete(name string) error {
	_, err := client.issueRequest("DELETE", name)

	if err != nil {
		return err
	}

	return nil
}

func (client *ExtractorsClient) GetAll() (ExtractorList, error) {

	response, err := client.issueRequest("GET", "")

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
