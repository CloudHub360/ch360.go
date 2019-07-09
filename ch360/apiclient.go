package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/CloudHub360/ch360.go/response"
	"io"
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

func NewApiClient(httpClient net.HttpDoer,
	apiUrl string,
	clientId string,
	clientSecret string,
	httpLogSink io.Writer) *ApiClient {

	var myHttpClient = httpClient

	if httpLogSink != nil {
		myHttpClient = NewLoggingDoer(myHttpClient, httpLogSink)
	}

	myHttpClient = net.NewUserAgentHttpClient(myHttpClient, "surf/"+Version)
	myHttpClient = net.NewContextAwareHttpClient(myHttpClient)

	tokenRetriever := NewTokenRetriever(myHttpClient, apiUrl)

	myHttpClient = &AuthorisingDoer{
		wrappedSender:  myHttpClient,
		tokenRetriever: tokenRetriever,
		clientId:       clientId,
		clientSecret:   clientSecret,
	}

	myHttpClient = &ResponseCheckingDoer{
		wrappedSender:   myHttpClient,
		responseChecker: &response.ErrorChecker{},
	}

	apiClient := &ApiClient{
		Classifiers: &ClassifiersClient{
			baseUrl:       apiUrl,
			requestSender: myHttpClient,
		},
		Documents: &DocumentsClient{
			baseUrl:       apiUrl,
			requestSender: myHttpClient,
		},
		Extractors: &ExtractorsClient{
			baseUrl:       apiUrl,
			requestSender: myHttpClient,
		},
		Modules: &ModulesClient{
			baseUrl:       apiUrl,
			requestSender: myHttpClient,
		},
	}

	return apiClient
}
