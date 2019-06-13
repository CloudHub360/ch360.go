package ch360

import (
	"context"
	"io"
	"net/http"
)

type requestBuilder struct {
	request *http.Request
	err     error
}

func newRequest(ctx context.Context, method string, url string, body io.Reader) *requestBuilder {
	request, err := http.NewRequest(method,
		url,
		body)

	if err != nil {
		return &requestBuilder{
			err: err,
		}
	}

	request = request.WithContext(ctx)

	return &requestBuilder{
		request: request,
	}
}

func (b *requestBuilder) withHeaders(headers map[string]string) *requestBuilder {
	if b.err != nil {
		return b
	}

	for k, v := range headers {
		b.request.Header.Add(k, v)
	}

	return b
}

func (b *requestBuilder) build() (*http.Request, error) {
	return b.request, b.err
}
