package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"os"
)

type Command interface {
	Execute(ctx context.Context) error
}

func CommandFor(runParams *config.RunParams, apiClient *ch360.ApiClient) (Command, error) {
	out := os.Stdout

	if runParams.Verb() == config.Classify {
		return NewClassifyFilesCommandFromArgs(runParams, apiClient)

	} else if runParams.Verb() == config.Extract {
		return NewExtractFilesCommandFromArgs(runParams, apiClient)

	} else if runParams.Verb() == config.Read {
		return NewReadFilesCommandFromArgs(runParams, apiClient)

	} else if runParams.Noun() == config.Classifier {
		switch runParams.Verb() {
		case config.Create:
			return NewCreateClassifierFromArgs(runParams, apiClient, out)
		case config.Upload:
			return NewUploadClassifierFromArgs(runParams, apiClient, out)
		case config.Delete:
			return NewDeleteClassifierFromArgs(runParams, apiClient.Classifiers, out)
		case config.List:
			return NewListClassifiers(apiClient.Classifiers, out), nil
		}

	} else if runParams.Noun() == config.Extractor {
		switch runParams.Verb() {
		case config.Create:
			if runParams.ModulesTemplate != "" {
				return NewCreateExtractorFromModulesWithArgs(runParams, apiClient.Extractors, out)
			}
			return NewCreateExtractorFromArgs(runParams, apiClient.Extractors, out)
		case config.Delete:
			return NewDeleteExtractorFromArgs(runParams, apiClient.Extractors, out)
		case config.List:
			return NewListExtractors(apiClient.Extractors, out), nil
		}
	} else if runParams.Noun() == config.Module {
		switch runParams.Verb() {
		case config.List:
			return NewListModules(apiClient.Modules, out), nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Unknown command: %s %s", runParams.Verb(), runParams.Noun()))
}
