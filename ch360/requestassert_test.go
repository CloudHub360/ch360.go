package ch360_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// helper type used to assert on various bits of the http request
type requestAssertion struct {
	request *http.Request
}

func (r requestAssertion) WithBody(t *testing.T, expectedBody []byte) requestAssertion {
	actualBody := bytes.Buffer{}
	_, _ = actualBody.ReadFrom(r.request.Body)
	assert.Equal(t, expectedBody, actualBody.Bytes())
	return r
}

func (r requestAssertion) WithHeaders(t *testing.T, headers map[string][]string) requestAssertion {

	for expectedHeader, expectedHeaderValue := range headers {

		actualHeaderValue, ok := r.request.Header[expectedHeader]

		assert.True(t, ok)
		assert.Equal(t, expectedHeaderValue, actualHeaderValue)
	}

	return r
}
