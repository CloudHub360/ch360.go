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
	Extractors  *ExtractorsClient
	Modules     *ModulesClient
}

func NewTokenRetriever(httpClient net.HttpDoer, apiUrl string) auth.TokenRetriever {
	return auth.NewHttpTokenCache(
		auth.NewHttpTokenRetriever(
			httpClient,
			apiUrl,
			&response.ErrorChecker{}))
}

func NewApiClient(httpClient net.HttpDoer, apiUrl string, clientId string, clientSecret string) *ApiClient {

	ctxhttpClient := net.NewContextAwareHttpClient(httpClient)

	tokenRetriever := NewTokenRetriever(ctxhttpClient, apiUrl)

	authorisingDoer := AuthorisingDoer{
		wrappedSender:  ctxhttpClient,
		tokenRetriever: tokenRetriever,
		clientId:       clientId,
		clientSecret:   clientSecret,
	}

	responseCheckingDoer := ResponseCheckingDoer{
		wrappedSender:   &authorisingDoer,
		responseChecker: &response.ErrorChecker{},
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
		Extractors: &ExtractorsClient{
			baseUrl:       apiUrl,
			requestSender: &responseCheckingDoer,
		},
		Modules: &ModulesClient{
			baseUrl:       apiUrl,
			requestSender: &responseCheckingDoer,
		},
	}

	return apiClient
}
