package authtoken

import (
	"net/http"
	"net/url"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
)

const ApiAddress = "https://api.cloudhub360.com"

type Getter interface {
	Get() (string, error) // todo custom type
}

type HttpGetter struct {
	apiUrl       string
	clientId     string
	clientSecret string
	httpClient   *http.Client
}

func NewHttpGetter(clientId string, clientSecret string, httpClient *http.Client, apiUrl string) *HttpGetter {
	return &HttpGetter{
		clientId:     clientId,
		httpClient:   httpClient,
		clientSecret: clientSecret,
		apiUrl:       apiUrl,
	}
}

func (getter *HttpGetter) Get() (string, error) {
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	type errorResponse struct {
		Message string `json:"message"`
	}

	form := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{getter.clientId},
		"client_secret": []string{getter.clientSecret},
	}

	resp, err := getter.httpClient.PostForm(getter.apiUrl+"/oauth/token", form)
	if err != nil {
		// No response received
		return "", err
	}

	defer resp.Body.Close()

	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)

	// Check status code
	if resp.StatusCode >= 400 {
		errResponse := errorResponse{}
		err = json.Unmarshal(buf.Bytes(), &errResponse)

		if err != nil || errResponse.Message == "" {
			return "", errors.New(fmt.Sprintf("An error occurred when requesting an authentication token (HTTP %d)", resp.StatusCode))
		}

		return "", errors.New(errResponse.Message)
	}

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf(
			"An unexpected response code was received when requesting an authentication token (HTTP %d)",
			resp.StatusCode))
	}

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
