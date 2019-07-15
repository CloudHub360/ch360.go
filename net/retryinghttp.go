package net

import (
	"bytes"
	"github.com/CloudHub360/ch360.go/ioutils"
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"io/ioutil"

	"net/http"
)

var _ HttpDoer = (*RetryingHttpClient)(nil)

// RetryingHttpClient is an HttpDoer decorator that retries HTTP requests
// for any response with HTTP status 500+, or any network error.
type RetryingHttpClient struct {
	wrapped       HttpDoer
	retryAttempts int
	multiplier    float64
}

func NewRetryingHttpClient(wrappedClient HttpDoer, retryAttempts int, multiplier float64) *RetryingHttpClient {
	return &RetryingHttpClient{
		wrapped:       wrappedClient,
		retryAttempts: retryAttempts,
		multiplier:    multiplier,
	}
}

func (h *RetryingHttpClient) Do(request *http.Request) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)
	requestBody, err := ioutils.DrainClose(request.Body)

	if err != nil {
		return nil, errors.WithMessage(err, "Unable to save request body")
	}

	var exponentialPolicy = backoff.NewExponentialBackOff()
	exponentialPolicy.Multiplier = h.multiplier
	backoffPolicy := backoff.WithMaxRetries(exponentialPolicy, uint64(h.retryAttempts))
	backoffPolicy = backoff.WithContext(backoffPolicy, request.Context())

	err = backoff.Retry(func() error {
		// Reset the body on the request to ensure it's readable (rewound)
		request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody.Bytes()))

		response, err = h.wrapped.Do(request)
		return shouldRetry(response, err)
	}, backoffPolicy)

	return response, err
}

// shouldRetry returns an error if the provided response and error are retryable.
func shouldRetry(response *http.Response, err error) error {
	if err != nil {
		return err
	}

	if response.StatusCode >= 500 || response.StatusCode == 408 {
		return errors.Errorf("Unexpected HTTP response %v", response.StatusCode)
	}

	// no error, don't retry
	return nil
}
