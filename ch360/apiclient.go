package ch360

import (
	"github.com/waives/surf/auth"
	"github.com/waives/surf/net"
	"io"
)

const ApiAddress = "https://api.cloudhub360.com"

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
			&net.ErrorChecker{}))
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
	myHttpClient = net.NewRetryingHttpClient(myHttpClient, 3, 2)

	tokenRetriever := NewTokenRetriever(myHttpClient, apiUrl)

	myHttpClient = &AuthorisingDoer{
		wrappedSender:  myHttpClient,
		tokenRetriever: tokenRetriever,
		clientId:       clientId,
		clientSecret:   clientSecret,
	}

	myHttpClient = &ResponseCheckingDoer{
		wrappedSender:   myHttpClient,
		responseChecker: &net.ErrorChecker{},
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
