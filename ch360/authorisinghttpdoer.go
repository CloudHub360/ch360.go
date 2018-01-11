package ch360

import (
	"github.com/CloudHub360/ch360.go/auth"
	"net/http"
)

type authorisingDoer struct {
	tokenRetriever auth.TokenRetriever
	wrappedSender  HttpDoer
}

func (sender *authorisingDoer) Do(request *http.Request) (*http.Response, error) {
	token, err := sender.tokenRetriever.RetrieveToken()

	if err != nil {
		return nil, err
	}

	if request.Header == nil {
		request.Header = make(http.Header)
	}

	request.Header.Add("Authorization", "Bearer "+token.TokenString)

	return sender.wrappedSender.Do(request)
}
