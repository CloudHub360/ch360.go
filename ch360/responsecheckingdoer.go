package ch360

import (
	"github.com/CloudHub360/ch360.go/response"
	"net/http"
)

type responseCheckingDoer struct {
	checker       response.Checker
	wrappedSender HttpDoer
}

func (sender *responseCheckingDoer) Do(request *http.Request) (*http.Response, error) {
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
