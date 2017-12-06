package auth

import (
	"bytes"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/response"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type TokenRetriever interface {
	RetrieveToken() (string, error)
}

type HttpTokenRetriever struct {
	apiUrl          string
	clientId        string
	clientSecret    string
	formPoster      FormPoster
	responseChecker response.Checker
}

type FormPoster interface {
	PostForm(url string, values url.Values) (*http.Response, error)
}

func NewHttpTokenRetriever(clientId string, clientSecret string, formPoster FormPoster, apiUrl string, responseChecker response.Checker) *HttpTokenRetriever {
	return &HttpTokenRetriever{
		clientId:        clientId,
		formPoster:      formPoster,
		clientSecret:    clientSecret,
		apiUrl:          apiUrl,
		responseChecker: responseChecker,
	}
}

func (retriever *HttpTokenRetriever) RetrieveToken() (string, error) {
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	form := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{retriever.clientId},
		"client_secret": []string{retriever.clientSecret},
	}

	resp, err := retriever.formPoster.PostForm(retriever.apiUrl+"/oauth/token", form)
	if err != nil {
		// No response received
		return "", err
	}
	defer resp.Body.Close()

	err = retriever.responseChecker.Check(resp)

	if err != nil {
		return "", errors.Wrap(err, "An error occurred when requesting an authentication token")
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)

	accessToken := tokenResponse{}
	err = json.Unmarshal(buf.Bytes(), &accessToken)

	if err != nil {
		return "", errors.New("Failed to parse authentication token response")
	}

	if accessToken.AccessToken == "" {
		return "", errors.New("Received empty authentication token")
	}
	return accessToken.AccessToken, nil
}
