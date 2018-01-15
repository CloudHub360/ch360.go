package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/response"
	"net/http"
)

const ApiAddress = "https://api.cloudhub360.com"

type ApiClient struct {
	Classifiers *ClassifiersClient
	Documents   *DocumentsClient
}

func NewApiClient(httpClient *http.Client, apiUrl string, clientId string, clientSecret string) *ApiClient {

	responseChecker := response.ErrorChecker{}

	tokenRetriever := auth.NewHttpTokenCache(
		auth.NewHttpTokenRetriever(
			clientId,
			clientSecret,
			httpClient,
			apiUrl,
			&responseChecker))

	authorisingDoer := AuthorisingDoer{
		wrappedSender:  httpClient,
		tokenRetriever: tokenRetriever,
	}

	responseCheckingDoer := ResponseCheckingDoer{
		wrappedSender:   &authorisingDoer,
		responseChecker: &responseChecker,
	}

	apiClient := &ApiClient{
		Classifiers: &ClassifiersClient{
			baseUrl:       apiUrl,
			requestSender: &responseCheckingDoer,
		},
		Documents: &DocumentsClient{
			baseUrl:       apiUrl,
			requestSender: &responseCheckingDoer,
		},
	}

	return apiClient
}
