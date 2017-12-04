package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var fakeClientId = "fake-client-id"
var fakeClientSecret = "fake-client-secret"

func Test_HttpTokenRetriever_Sends_Client_Id_And_Secret(t *testing.T) {
	// Arrange
	var receivedClientId string
	var receivedClientSecret string

	requestHandler := func(w http.ResponseWriter, r *http.Request) {
		receivedClientId = r.FormValue("client_id")
		receivedClientSecret = r.FormValue("client_secret")
		w.WriteHeader(200)
	}

	// create test server with requestHandler
	server := httptest.NewServer(http.HandlerFunc(requestHandler))
	defer server.Close()

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, server.Client(), server.URL)

	// Act
	sut.RetrieveToken()

	// Assert
	if receivedClientId != fakeClientId {
		t.Error("Did not receive client ID")
	}

	if receivedClientSecret != fakeClientSecret {
		t.Error("Did not receive client secret")
	}
}

func Test_HttpTokenRetriever_Parses_Json_Response(t *testing.T) {
	// Arrange
	fakeServer := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token": "%s"}`, "fake-token")
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(fakeServer))
	defer ts.Close()

	sut := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, ts.Client(), ts.URL)

	// Act
	token, err := sut.RetrieveToken()

	// Assert
	if err != nil {
		t.Error(err)
	}

	if token != "fake-token" {
		t.Error("Incorrect auth token")
	}
}

var unsuccessfulRequestData = []struct {
	responseCode int
	responseBody []byte
	expectedErr  string
}{
	{201, nil, "An error occurred when requesting an authentication token: Received unexpected response code: 201"},
	{200, []byte(`{"access_token": ""}`), "Received empty authentication token"},
	{200, []byte(`<invalid json>`), "Failed to parse authentication token response"},
	{400, []byte(`{"message": "error-message"}`), "An error occurred when requesting an authentication token: error-message"},
	{499, []byte(`{"message": "error-message"}`), "An error occurred when requesting an authentication token: error-message"},
	{403, []byte(`<Invalid json>`), "An error occurred when requesting an authentication token: Received error response with HTTP code 403"},
	{500, nil, "An error occurred when requesting an authentication token: Received error response with HTTP code 500"},
	{501, nil, "An error occurred when requesting an authentication token: Received error response with HTTP code 501"},
}

func Test_HttpTokenRetriever_Returns_Err_On_Unsuccessful_Request(t *testing.T) {
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

			tokenGetter := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, ts.Client(), ts.URL)

			// Act
			_, err := tokenGetter.RetrieveToken()

			// Assert
			if err == nil {
				t.Error("HttpTokenRetriever.RetrieveToken() didn't return an error, but should have.", tp.expectedErr)
			} else if err.Error() != tp.expectedErr {
				t.Error("Incorrect error message received", tp.expectedErr, err.Error())
			}
		}()
	}
}

func Test_HttpTokenRetriever_Returns_Err_On_Client_Error(t *testing.T) {
	// Arrange
	tokenGetter := NewHttpTokenRetriever(fakeClientId, fakeClientSecret, &http.Client{}, "http://invalid-url:-1")

	// Act
	_, err := tokenGetter.RetrieveToken()

	// Assert
	if err == nil {
		t.Error("HttpTokenRetriever.RetrieveToken() didn't return an error, but should have.")
	}
}
