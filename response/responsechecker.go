package response

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrorChecker struct{}

type Checker interface {
	CheckForErrors(response *http.Response) error
}

func (c *ErrorChecker) CheckForErrors(response *http.Response) error {
	type errorResponse struct {
		Message string `json:"message"`
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(response.Body)

	// We've read from the response body, and it can't be rewound, so 'recreate' it as a new io.Reader
	// which will read from the start of the underlying byte array of 'buf'.
	response.Body = ioutil.NopCloser(bufio.NewReader(&buf))

	// Check status code
	if response.StatusCode < 400 {
		return nil
	}

	errResponse := errorResponse{}
	err := json.Unmarshal(buf.Bytes(), &errResponse)

	if err != nil || errResponse.Message == "" {
		return errors.New(fmt.Sprintf("Received error response with HTTP code %d", response.StatusCode))
	}

	return errors.New(errResponse.Message)

}
