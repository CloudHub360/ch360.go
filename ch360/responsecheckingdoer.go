package ch360

import (
	"github.com/CloudHub360/ch360.go/net"
	"github.com/CloudHub360/ch360.go/response"
	"net/http"
)

type ResponseCheckingDoer struct {
	responseChecker response.Checker
	wrappedSender   net.HttpDoer
}

func NewResponseCheckingdoer(checker response.Checker, wrappedSender net.HttpDoer) *ResponseCheckingDoer {
	return &ResponseCheckingDoer{
		wrappedSender:   wrappedSender,
		responseChecker: checker,
	}
}

func (requestSender *ResponseCheckingDoer) Do(request *http.Request) (*http.Response, error) {
	response, err := requestSender.wrappedSender.Do(request)

	if err != nil {
		return nil, err
	}

	err = requestSender.responseChecker.CheckForErrors(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
