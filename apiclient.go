package ch360

import (
	"io"
	"net/http"
	"github.com/CloudHub360/ch360.go/authtoken"
	"github.com/CloudHub360/ch360.go/response"
)

const ApiAddress = "https://api.cloudhub360.com"

type ApiClient struct {
	apiUrl       string
	tokenGetter  authtoken.Getter
	httpClient   *http.Client
	responseChecker response.Checker
}

func NewApiClient(httpClient *http.Client, apiUrl string, tokenGetter authtoken.Getter) (*ApiClient) {
	apiClient := &ApiClient{
		apiUrl:       apiUrl,
		httpClient:   httpClient,
		tokenGetter:  tokenGetter,
		responseChecker:response.Checker{},
	}

	return apiClient
}

func (ac *ApiClient) send(method string, path string, body io.Reader) ([]byte, error) {
	token, err := ac.tokenGetter.Get()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, ac.apiUrl+path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+ token)

	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ac.responseChecker.Check(resp, 200)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (ac *ApiClient) CreateClassifier(name string) (error) {
	_, err := ac.send("POST", "/classifiers/" + name, nil)
	return err
}