package net

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func Test_Returns_Error_With_Correct_Message_When_Check_Fails(t *testing.T) {
	var fixtures = []struct {
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

	for _, fixture := range fixtures {
		sut := &ErrorChecker{}

		response := http.Response{
			StatusCode: fixture.responseCode,
			Body:       ioutil.NopCloser(bytes.NewBuffer(fixture.responseBody)),
		}

		err := sut.CheckForErrors(&response)

		assert.EqualError(t, err, fixture.expectedErr)
	}
}

func Test_Returns_Correct_Error_Type_For_Rfc7807_And_Simple_Responses(t *testing.T) {
	var fixtures = []struct {
		responseCode      int
		responseBody      []byte
		expectedErrorType reflect.Type
	}{
		{422, []byte(rfc7807Response), reflect.TypeOf(&DetailedErrorResponse{})},
		{422, []byte(`{"message": "error-message"}`), reflect.TypeOf(&basicErrorResponse{})},
	}

	for _, fixture := range fixtures {
		// Arrange
		sut := &ErrorChecker{}
		response := http.Response{
			StatusCode: fixture.responseCode,
			Body:       ioutil.NopCloser(bytes.NewBuffer(fixture.responseBody)),
		}

		// Act
		err := sut.CheckForErrors(&response)

		// Assert
		assert.Equal(t, fixture.expectedErrorType, reflect.TypeOf(err))
	}
}

var rfc7807Response = `{
  "errors": [
    {
      "module_id": "waives.supplier_identity",
      "messages": [
        "No argument was specified"
      ],
      "path": "modules[0].arguments.provider",
      "argument_name": "provider",
      "argument_value": ""
    }
  ],
  "type": "https://docs.waives.io/reference#invalid-extractor-template",
  "title": "Invalid Extractor Template",
  "status": 422,
  "instance": "/account/jK16_1URgUGxSo6yyWjHag/invalid-extractor-template/supplier-id"
}`
