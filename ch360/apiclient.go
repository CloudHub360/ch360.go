package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/response"
	"io"
	"net/http"
)

const ApiAddress = "https://api.cloudhub360.com"

type ApiClient struct {
	apiUrl          string
	retriever       auth.TokenRetriever
	httpClient      DoRequester
	responseChecker response.ErrorChecker
	Classifiers     *ClassifiersClient
}

//TODO: Should this take a pointer to a requester (which is an httpClient)?
func NewApiClient(requester DoRequester, apiUrl string, tokenRetriever auth.TokenRetriever) *ApiClient {
	apiClient := &ApiClient{
		apiUrl:          apiUrl,
		httpClient:      requester,
		retriever:       tokenRetriever,
		responseChecker: response.ErrorChecker{},
		Classifiers:     &ClassifiersClient{},
	}

	return apiClient
}

type Sender interface {
	Send(method string, path string, body io.Reader) ([]byte, error)
}

type DoRequester interface {
	Do(request *http.Request) (*http.Response, error)
}

func (client *ApiClient) Send(method string, path string, body io.Reader) ([]byte, error) {
	token, err := client.retriever.RetrieveToken()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, client.apiUrl+path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := client.responseChecker.Check(resp, 200)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
