package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"net/http"
	"os"
	"time"
)

type Command interface {
	Execute(ctx context.Context) error
}

func CommandFor(runParams *config.RunParams) (Command, error) {
	out := os.Stdout

	apiClient, err := initSurf(runParams)

	if err != nil {
		return nil, err
	}

	if runParams.Verb() == config.Classify {
		return NewClassifyFilesCommandFromArgs(runParams, apiClient)

	} else if runParams.Verb() == config.Extract {
		return NewExtractFilesCommandFromArgs(runParams, apiClient)

	} else if runParams.Noun() == config.Classifier {
		switch runParams.Verb() {
		case config.Create:
			return NewCreateClassifierFromArgs(runParams, apiClient, out)
		case config.Delete:
			return NewDeleteClassifierFromArgs(runParams, apiClient.Classifiers, out)
		case config.List:
			return NewListClassifiers(apiClient.Classifiers, out), nil
		}

	} else if runParams.Noun() == config.Extractor {
		switch runParams.Verb() {
		case config.Create:
			return NewCreateExtractorFromArgs(runParams, apiClient.Extractors, out)
		case config.Delete:
			return NewDeleteExtractorFromArgs(runParams, apiClient.Extractors, out)
		case config.List:
			return NewListExtractors(apiClient.Extractors, out), nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Unknown command: %s %s", runParams.Verb(), runParams.Noun()))
}

var DefaultHttpClient = &http.Client{Timeout: time.Minute * 5}

func initSurf(params *config.RunParams) (*ch360.ApiClient, error) {

	appDir, err := config.NewAppDirectory()
	if err != nil {
		return nil, err
	}

	credentialsResolver := &CredentialsResolver{}

	clientId, clientSecret, err := credentialsResolver.ResolveFromArgs(params.Args(), appDir)

	if err != nil {
		return nil, err
	}

	return ch360.NewApiClient(DefaultHttpClient, ch360.ApiAddress, clientId, clientSecret), nil
}
