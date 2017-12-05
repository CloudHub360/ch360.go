package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/response"
	"net/http"
)

const ApiAddress = "https://api.cloudhub360.com"

type ApiClient struct {
	Classifiers *ClassifiersClient
}

func NewApiClient(httpClient *http.Client, apiUrl string, clientId string, clientSecret string) *ApiClient {

	responseChecker := response.ErrorChecker{}

	tokenRetriever := auth.NewHttpTokenRetriever(clientId,
		clientSecret, httpClient, apiUrl, &responseChecker)

	authorisingSender := AuthorisingSender{
		wrappedSender:httpClient,
		retriever:tokenRetriever,
	}

	responseCheckingSender := ResponseCheckingSender{
		wrappedSender:&authorisingSender,
		checker:&responseChecker,
	}

	apiClient := &ApiClient{
		Classifiers: &ClassifiersClient{
			baseUrl: apiUrl,
			sender:  &responseCheckingSender,
		},
	}

	return apiClient
}

type HttpSender interface {
	Do(request *http.Request) (*http.Response, error)
}

type ResponseCheckingSender struct {
	checker       response.Checker
	wrappedSender HttpSender
}

func (sender *ResponseCheckingSender) Do(request *http.Request) (*http.Response, error) {
	response, err := sender.wrappedSender.Do(request)

	err = sender.checker.Check(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type AuthorisingSender struct {
	retriever     auth.TokenRetriever
	wrappedSender HttpSender
}

func (sender *AuthorisingSender) Do(request *http.Request) (*http.Response, error) {
	token, err := sender.retriever.RetrieveToken()

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token)

	return sender.wrappedSender.Do(request)
}
