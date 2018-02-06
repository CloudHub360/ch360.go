package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/net"
	"net/http"
)

type AuthorisingDoer struct {
	tokenRetriever auth.TokenRetriever
	wrappedSender  net.HttpDoer
	clientId       string
	clientSecret   string
}

func NewAuthorisingDoer(retriever auth.TokenRetriever, httpDoer net.HttpDoer, clientId string, clientSecret string) *AuthorisingDoer {
	return &AuthorisingDoer{
		tokenRetriever: retriever,
		wrappedSender:  httpDoer,
		clientId:       clientId,
		clientSecret:   clientSecret,
	}
}

func (ad *AuthorisingDoer) Do(request *http.Request) (*http.Response, error) {
	token, err := ad.tokenRetriever.RetrieveToken(ad.clientId, ad.clientSecret)

	if err != nil {
		return nil, err
	}

	if request.Header == nil {
		request.Header = make(http.Header)
	}

	request.Header.Add("Authorization", "Bearer "+token.TokenString)

	return ad.wrappedSender.Do(request)
}
