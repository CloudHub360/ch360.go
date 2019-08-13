package commands

import (
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"io"
	"net/http"
	"os"
	"time"
)

func initApiClient(clientIdFlag, clientSecretFlag string, logHttpFile *os.File) (*ch360.ApiClient, error) {
	appDir, err := config.NewAppDirectory()
	if err != nil {
		return nil, err
	}

	credentialsResolver := &CredentialsResolver{}

	clientId, clientSecret, err := credentialsResolver.Resolve(clientIdFlag, clientSecretFlag, appDir)

	if err != nil {
		return nil, err
	}

	var logSink io.Writer = nil
	if logHttpFile != nil {
		logSink = logHttpFile
	}
	return ch360.NewApiClient(DefaultHttpClient, ch360.ApiAddress, clientId, clientSecret, logSink), nil
}

var DefaultHttpClient = &http.Client{Timeout: time.Minute * 2}
