package ch360

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"github.com/CloudHub360/ch360.go/authtoken"
)

type ApiClient struct {
	apiUrl       string
	tokenGetter  authtoken.Getter
	httpClient   *http.Client
}

func NewApiClient(httpClient *http.Client, apiUrl string, tokenGetter authtoken.Getter) (*ApiClient) {
	apiClient := &ApiClient{
		apiUrl:       apiUrl,
		httpClient:   httpClient,
		tokenGetter:  tokenGetter,
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

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("Request failed with code %d: [%s] %s", resp.StatusCode, method, path))
	}

	defer resp.Body.Close()
	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

func (ac *ApiClient) CreateClassifier(name string) (error) {
	_, err := ac.send("POST", "/classifiers/" + name, nil)
	return err
}