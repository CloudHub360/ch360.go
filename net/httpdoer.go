package net

import (
	"net/http"
)

//go:generate mockery -name HttpDoer

type HttpDoer interface {
	Do(request *http.Request) (*http.Response, error)
}
