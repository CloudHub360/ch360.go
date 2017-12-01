package authtoken

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"net/http"
)

var fakeClientId = "fake-client-id"
var fakeClientSecret = "fake-client-secret"

func Test_HttpGetter_Sends_Client_Id_And_Secret(t *testing.T) {
	// Arrange
	var receivedClientId string
	var receivedClientSecret string

	fakeServer := func(w http.ResponseWriter, r *http.Request) {
		receivedClientId = r.FormValue("client_id")
		receivedClientSecret = r.FormValue("client_secret")
		w.WriteHeader(200)
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(fakeServer))
	defer ts.Close()

	sut := NewHttpGetter(fakeClientId, fakeClientSecret, ts.Client(), ts.URL);

	// Act
	sut.Get()

	// Assert
	if receivedClientId != fakeClientId {
		t.Error("Did not receive client ID")
	}

	if receivedClientSecret != fakeClientSecret {
		t.Error("Did not receive client secret")
	}
}

func Test_HttpGetter_Parses_Json_Response(t *testing.T) {
	// Arrange
	fakeServer := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token": "%s"}`, "fake-token")
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(fakeServer))
	defer ts.Close()

	sut := NewHttpGetter(fakeClientId, fakeClientSecret, ts.Client(), ts.URL);

	// Act
	token, err := sut.Get()

	// Assert
	if err != nil {
		t.Error(err)
	}

	if token != "fake-token" {
		t.Error("Incorrect auth token")
	}
}

var unsuccessfulRequestData = []struct {
	responseCode  int
	responseBody    []byte
	expectedErr		string
}{
	{201, nil, "An unexpected response code was received when " +
		"requesting an authentication token (HTTP 201)"},
	{200, []byte(`{"access_token": ""}`), "Received empty authentication token"},
	{200, []byte(`<invalid json>`), "Failed to parse authentication token response"},
	{400, []byte(`{"message": "error-message"}`), "error-message"},
	{499, []byte(`{"message": "error-message"}`), "error-message"},
	{403, []byte(`<Invalid json>`), "An error occurred when requesting an " +
		"authentication token (HTTP 403)"},
	{500, nil, "An error occurred when requesting an authentication token (HTTP 500)"},
	{501, nil, "An error occurred when requesting an authentication token (HTTP 501)"},
}
func Test_HttpGetter_Returns_Err_On_Unsuccessful_Request(t *testing.T) {
	for _, tp := range unsuccessfulRequestData {
		// run an anonymous function to ensure defer is called on each iteration
		func() {
			// Arrange
			fakeServer := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tp.responseCode)
				w.Write(tp.responseBody)
			}

			ts := httptest.NewServer(http.HandlerFunc(fakeServer))
			defer ts.Close()

			tokenGetter := NewHttpGetter(fakeClientId, fakeClientSecret, ts.Client(), ts.URL)

			// Act
			_, err := tokenGetter.Get()

			// Assert
			if err == nil {
				t.Error("HttpGetter.Get() didn't return an error, but should have.", tp.expectedErr)
			} else if err.Error() != tp.expectedErr {
				t.Error("Incorrect error message received" , tp.expectedErr, err.Error())
			}
		}()
	}
}

func Test_HttpGetter_Returns_Err_On_Client_Error(t *testing.T) {
	// Arrange
	tokenGetter := NewHttpGetter(fakeClientId, fakeClientSecret, &http.Client{}, "http://invalid-url:-1")

	// Act
	_, err := tokenGetter.Get()

	// Assert
	if err == nil {
		t.Error("HttpGetter.Get() didn't return an error, but should have.")
	}
}
