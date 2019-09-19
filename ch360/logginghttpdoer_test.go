package ch360

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/waives/surf/net/mocks"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_LoggingHttpDoer_Logs_To_Provided_Sink(t *testing.T) {
	fixtures := []struct {
		requestBody    []byte
		responseBody   []byte
		isRequestJson  bool
		isResponseJson bool
	}{
		{
			requestBody:    nil,
			isRequestJson:  false,
			responseBody:   []byte("non-json-body"),
			isResponseJson: false,
		}, {
			requestBody:    []byte("some binary bytes go here"),
			isRequestJson:  false,
			responseBody:   []byte(`{"jsonmessage":"all good"}`),
			isResponseJson: true, // json response body should be logged
		}, {
			requestBody:    []byte(`{"jsonmessage":"json request body"}`),
			isRequestJson:  true, // json request body should be logged
			responseBody:   []byte("non-json-body"),
			isResponseJson: false,
		},
	}

	for _, fixture := range fixtures {

		// Arrange
		httpDoer := mocks.HttpDoer{}
		logSink := bytes.Buffer{}
		sut := LoggingDoer{
			wrappedSender: &httpDoer,
			out:           &logSink,
		}
		request, _ := http.NewRequest("GET", "https://api.waives.io", bytes.NewBuffer(fixture.requestBody))
		response := http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBuffer(fixture.responseBody)),
		}
		httpDoer.
			On("Do", request).
			Return(&response, nil)

		// Act
		_, _ = sut.Do(request)

		// Assert
		assert.Contains(t, logSink.String(), "GET / HTTP/1.1\r\nHost: api.waives.io")
		if fixture.isRequestJson {
			assert.Contains(t, logSink.String(), formatJson(string(fixture.requestBody)))
		}
		if fixture.isResponseJson {
			assert.Contains(t, logSink.String(), formatJson(string(fixture.responseBody)))
		}
	}
}

func formatJson(input string) string {
	dst := bytes.Buffer{}
	err := json.Indent(&dst, []byte(input), "", "  ")

	if err != nil {
		panic(err)
	}

	return dst.String()
}
