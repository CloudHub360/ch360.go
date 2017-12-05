package ch360

import (
	"github.com/CloudHub360/ch360.go/response"
	"net/http"
	"github.com/CloudHub360/ch360.go/auth"
)

type HttpDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

type ResponseCheckingDoer struct {
	checker       response.Checker
	wrappedSender HttpDoer
}

func (sender *ResponseCheckingDoer) Do(request *http.Request) (*http.Response, error) {
	response, err := sender.wrappedSender.Do(request)

	if err != nil {
		return nil, err
	}

	err = sender.checker.Check(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type AuthorisingDoer struct {
	retriever     auth.TokenRetriever
	wrappedSender HttpDoer
}

func (sender *AuthorisingDoer) Do(request *http.Request) (*http.Response, error) {
	token, err := sender.retriever.RetrieveToken()

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token)

	return sender.wrappedSender.Do(request)
}
