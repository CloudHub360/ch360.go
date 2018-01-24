package net

import (
	"net/http"
)

type ContextAwareHttpClient struct {
	httpDoer HttpDoer
}

func (c *ContextAwareHttpClient) Do(r *http.Request) (*http.Response, error) {
	ctx := r.Context()

	resp, err := c.httpDoer.Do(r)

	select {
	case <-ctx.Done():
		// If the context was cancelled, use its err message as it's more likely to be useful...
		err = ctx.Err()
	default:
		// ...otherwise, don't.
	}

	return resp, err
}

func NewContextAwareHttpClient(doer HttpDoer) *ContextAwareHttpClient {
	return &ContextAwareHttpClient{httpDoer: doer}
}
