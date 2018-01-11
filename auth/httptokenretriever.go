package auth

import (
	"bytes"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/response"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"time"
)

//go:generate mockery -name "TokenRetriever|FormPoster"

type TokenRetriever interface {
	RetrieveToken() (*AccessToken, error)
}

type AccessToken struct {
	TokenString string
	ExpiresAt   time.Time
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

func (retriever *HttpTokenRetriever) RetrieveToken() (*AccessToken, error) {
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	form := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{retriever.clientId},
		"client_secret": []string{retriever.clientSecret},
	}

	resp, err := retriever.formPoster.PostForm(retriever.apiUrl+"/oauth/token", form)
	if err != nil {
		// No response received
		return nil, err
	}
	defer resp.Body.Close()

	err = retriever.responseChecker.CheckForErrors(resp)

	if err != nil {
		return nil, errors.Wrap(err, "An error occurred when requesting an authentication token")
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)

	accessToken := tokenResponse{}
	err = json.Unmarshal(buf.Bytes(), &accessToken)

	if err != nil {
		return nil, errors.New("Failed to parse authentication token response")
	}

	if accessToken.AccessToken == "" {
		return nil, errors.New("Received empty authentication token")
	}
	return &AccessToken{
		accessToken.AccessToken,
		time.Now().Add(time.Duration(accessToken.ExpiresIn) * time.Minute),
	}, nil
}
