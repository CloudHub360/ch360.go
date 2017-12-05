package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"io/ioutil"
	"bytes"
	"errors"
)

func Test_Returns_Response_Body_When_Check_Passes(t *testing.T) {
	expectedResponseBody := []byte("body")
	sut := &ErrorChecker{}

	response := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(expectedResponseBody)),
	}

	body, err := sut.Check(&response, response.StatusCode)

	assert.Equal(t, expectedResponseBody, body)
	assert.Nil(t, err)
}

func Test_Returns_Error_With_Correct_Message_When_Check_Fails(t *testing.T) {
	for _, tp := range errorResponsesData {
		// run an anonymous function to ensure defer is called on each iteration
		func() {
			sut := &ErrorChecker{}

			response := http.Response{
				StatusCode: tp.responseCode,
				Body:       ioutil.NopCloser(bytes.NewBuffer(tp.responseBody)),
			}

			body, err := sut.Check(&response, 200)

			assert.Equal(t, errors.New(tp.expectedErr), err)
			assert.Nil(t, body)
		}()
	}
}

var errorResponsesData = []struct {
	responseCode int
	responseBody []byte
	expectedErr  string
}{
	{301, []byte(`<invalid json>`), "Received unexpected response code: 301"},
	{400, []byte(`{"message": "error-message"}`), "error-message"},
	{499, []byte(`{"message": "error-message"}`), "error-message"},
	{403, []byte(`<Invalid json>`), "Received error response with HTTP code 403"},
	{500, nil, "Received error response with HTTP code 500"},
	{501, nil, "Received error response with HTTP code 501"},
}