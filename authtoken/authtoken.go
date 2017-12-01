package authtoken

import (
	"net/http"
	"net/url"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/CloudHub360/ch360.go/response"
)


type Getter interface {
	Get() (string, error) // todo custom type
}

type HttpGetter struct {
	apiUrl       string
	clientId     string
	clientSecret string
	httpClient   *http.Client
	responseChecker response.Checker
}

func NewHttpGetter(clientId string, clientSecret string, httpClient *http.Client, apiUrl string) *HttpGetter {
	return &HttpGetter{
		clientId:     clientId,
		httpClient:   httpClient,
		clientSecret: clientSecret,
		apiUrl:       apiUrl,
		responseChecker: response.Checker{},
	}
}

func (getter *HttpGetter) Get() (string, error) {
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
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

	bytes, err := getter.responseChecker.Check(resp, 200)

	if err != nil {
		return "", errors.Wrap(err, "An error occurred when requesting an authentication token")
	}

	accessToken := tokenResponse{}
	err = json.Unmarshal(bytes, &accessToken)

	if err != nil {
		return "", errors.New("Failed to parse authentication token response")
	}

	if accessToken.AccessToken == "" {
		return "", errors.New("Received empty authentication token")
	}
	return accessToken.AccessToken, nil
}
