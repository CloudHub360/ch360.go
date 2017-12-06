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

	authorisingDoer := authorisingDoer{
		wrappedSender:  httpClient,
		tokenRetriever: tokenRetriever,
	}

	responseCheckingDoer := responseCheckingDoer{
		wrappedSender: &authorisingDoer,
		checker:       &responseChecker,
	}

	apiClient := &ApiClient{
		Classifiers: &ClassifiersClient{
			baseUrl: apiUrl,
			sender:  &responseCheckingDoer,
		},
	}

	return apiClient
}
