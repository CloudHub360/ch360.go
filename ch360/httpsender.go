package ch360

import (
	"github.com/CloudHub360/ch360.go/response"
	"net/http"
	"github.com/CloudHub360/ch360.go/auth"
)

type HttpSender interface {
	Do(request *http.Request) (*http.Response, error)
}

type ResponseCheckingSender struct {
	checker       response.Checker
	wrappedSender HttpSender
}

func (sender *ResponseCheckingSender) Do(request *http.Request) (*http.Response, error) {
	response, err := sender.wrappedSender.Do(request)

	err = sender.checker.Check(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type AuthorisingSender struct {
	retriever     auth.TokenRetriever
	wrappedSender HttpSender
}

func (sender *AuthorisingSender) Do(request *http.Request) (*http.Response, error) {
	token, err := sender.retriever.RetrieveToken()

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token)

	return sender.wrappedSender.Do(request)
}
