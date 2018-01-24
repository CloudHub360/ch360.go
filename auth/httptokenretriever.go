package auth

import (
	"bytes"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/CloudHub360/ch360.go/response"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
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
	httpDoer        net.HttpDoer
	responseChecker response.Checker
}

func NewHttpTokenRetriever(clientId string, clientSecret string, httpDoer net.HttpDoer, apiUrl string, responseChecker response.Checker) *HttpTokenRetriever {
	return &HttpTokenRetriever{
		clientId:        clientId,
		httpDoer:        httpDoer,
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

	req, err := http.NewRequest("POST", retriever.apiUrl+"/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := retriever.httpDoer.Do(req)
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
