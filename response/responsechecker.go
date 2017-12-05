package response

import (
	"net/http"
	"encoding/json"
	"fmt"
	"errors"
	"bytes"
)

type ErrorChecker struct {}

type Checker interface {
	Check(response *http.Response, successCode int) ([]byte, error)
}

func (c *ErrorChecker) Check(response *http.Response, successCode int) ([]byte, error) {
	type errorResponse struct {
		Message string `json:"message"`
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(response.Body)

	// Check status code
	if response.StatusCode >= 400 {
		errResponse := errorResponse{}
		err := json.Unmarshal(buf.Bytes(), &errResponse)

		if err != nil || errResponse.Message == "" {
			return nil, errors.New(fmt.Sprintf("Received error response with HTTP code %d", response.StatusCode))
		}

		return nil, errors.New(errResponse.Message)
	}

	if response.StatusCode != successCode {
		return nil, errors.New(fmt.Sprintf(
			"Received unexpected response code: %d", response.StatusCode))
	}

	return buf.Bytes(), nil
}