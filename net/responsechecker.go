package net

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type ErrorChecker struct{}

//go:generate mockery -name ResponseChecker

var _ ResponseChecker = (*ErrorChecker)(nil)

type ResponseChecker interface {
	CheckForErrors(response *http.Response) error
}

func (c *ErrorChecker) CheckForErrors(response *http.Response) error {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(response.Body)

	if err != nil {
		return errors.WithMessage(err, "Unable to read from HTTP response body")
	}
	response.Body.Close()

	// We've read from the response body, and it can't be rewound, so 'recreate' it as a new io.Reader
	// which will read from the start of the underlying byte array of 'buf'.
	response.Body = ioutil.NopCloser(&buf)

	// Check status code
	if response.StatusCode < 300 {
		return nil
	}

	if json.Valid(buf.Bytes()) {
		var (
			basicError    = &basicErrorResponse{}
			detailedError = &DetailedErrorResponse{}
		)
		// Try the basic err json first...
		err = json.Unmarshal(buf.Bytes(), &basicError)

		if err == nil && len(basicError.Message) > 0 {
			return basicError
		}

		// .. then the more detailed form
		err = json.Unmarshal(buf.Bytes(), &detailedError)

		if err == nil && detailedError.Status != 0 {
			return detailedError
		}
	}

	return errors.Errorf("Received unexpected response with HTTP code %d", response.StatusCode)
}

type basicErrorResponse struct {
	Message string `json:"message"`
}

func (e *basicErrorResponse) Error() string {
	return e.Message
}

type DetailedErrorResponse struct {
	Errors   []map[string]interface{} `json:"errors"`
	Type     string                   `json:"type"`
	Title    string                   `json:"title"`
	Status   int                      `json:"status"`
	Instance string                   `json:"instance"`
	Detail   string                   `json:"detail"`
}

func (e *DetailedErrorResponse) Error() string {
	if len(e.Detail) > 0 {
		return e.Detail
	}
	return e.Title
}
