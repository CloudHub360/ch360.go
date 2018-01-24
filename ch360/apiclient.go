package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/CloudHub360/ch360.go/response"
)

const ApiAddress = "https://api.waives.io"

type ApiClient struct {
	Classifiers *ClassifiersClient
	Documents   *DocumentsClient
}

func NewApiClient(httpClient net.HttpDoer, apiUrl string, clientId string, clientSecret string) *ApiClient {

	responseChecker := response.ErrorChecker{}

	ctxhttpClient := net.NewContextAwareHttpClient(httpClient)

	tokenRetriever := auth.NewHttpTokenCache(
		auth.NewHttpTokenRetriever(
			clientId,
			clientSecret,
			ctxhttpClient,
			apiUrl,
			&responseChecker))

	authorisingDoer := AuthorisingDoer{
		wrappedSender:  ctxhttpClient,
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
