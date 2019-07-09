package net

import "net/http"

var _ HttpDoer = (*UserAgentHttpClient)(nil)

// UserAgentHttpClient is an HttpDoer decorator that sets the "User-Agent"
// HTTP header on requests.
type UserAgentHttpClient struct {
	wrapper   HttpDoer
	userAgent string
}

func (h *UserAgentHttpClient) Do(request *http.Request) (*http.Response, error) {
	if request != nil {
		request.Header.Add("User-Agent", h.userAgent)
	}

	return h.wrapper.Do(request)
}

func NewUserAgentHttpClient(wrapper HttpDoer, userAgent string) *UserAgentHttpClient {
	return &UserAgentHttpClient{
		wrapper:   wrapper,
		userAgent: userAgent,
	}
}
