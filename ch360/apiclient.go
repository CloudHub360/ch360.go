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

	authorisingSender := authorisingDoer{
		wrappedSender:httpClient,
		retriever:tokenRetriever,
	}

	responseCheckingSender := responseCheckingDoer{
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

