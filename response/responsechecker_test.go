package response

import (
	"bytes"
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

			err := sut.CheckForErrors(&response)

			assert.EqualError(t, err, errorResponseData.expectedErr)
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
	{403, []byte(`<Invalid json>`), "Received unexpected response with HTTP code 403"},
	{400, nil, "Received unexpected response with HTTP code 400"},
	{500, nil, "Received unexpected response with HTTP code 500"},
	{501, nil, "Received unexpected response with HTTP code 501"},
	{502, []byte(`Bad gateway`), "Received unexpected response with HTTP code 502"},
	{300, nil, "Received unexpected response with HTTP code 300"},
	{301, []byte(`Moved permanently`), "Received unexpected response with HTTP code 301"},
}
