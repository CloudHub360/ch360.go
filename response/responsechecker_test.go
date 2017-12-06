package response

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Returns_Error_With_Correct_Message_When_Check_Fails(t *testing.T) {
	for _, errorResponseData := range errorResponsesData {
		// run an anonymous function to ensure defer is called on each iteration
		func() {
			sut := &ErrorChecker{}

			response := http.Response{
				StatusCode: errorResponseData.responseCode,
				Body:       ioutil.NopCloser(bytes.NewBuffer(errorResponseData.responseBody)),
			}

			err := sut.Check(&response)

			assert.Equal(t, errors.New(errorResponseData.expectedErr), err)
		}()
	}
}

var errorResponsesData = []struct {
	responseCode int
	responseBody []byte
	expectedErr  string
}{
	{400, []byte(`{"message": "error-message"}`), "error-message"},
	{499, []byte(`{"message": "error-message"}`), "error-message"},
	{403, []byte(`<Invalid json>`), "Received error response with HTTP code 403"},
	{500, nil, "Received error response with HTTP code 500"},
	{501, nil, "Received error response with HTTP code 501"},
}
