package ch360

import (
	"net/http"
)

type HttpDoer interface {
	Do(request *http.Request) (*http.Response, error)
}
