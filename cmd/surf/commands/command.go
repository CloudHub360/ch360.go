package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"time"
)

type Command interface {
	Execute(ctx context.Context) error
}

func CommandFor(args docopt.Opts) (Command, error) {
	cmd, err := parseCommand(args)

	if err != nil {
		return nil, err
	}

	apiClient, err := initSurf(args)

	if err != nil {
		return nil, err
	}

	out := os.Stdout

	switch cmd {
	case new(CreateClassifier).Usage():
		return NewCreateClassifierFromArgs(args, apiClient, out)
	case new(CreateExtractor).Usage():
		return NewCreateExtractorFromArgs(args, apiClient.Extractors, out)
	case new(DeleteClassifier).Usage():
		return NewDeleteClassifierFromArgs(args, apiClient.Classifiers, out)
	case new(ListClassifiers).Usage():
		return NewListClassifiers(apiClient.Classifiers, out), nil
	case new(ListExtractors).Usage():
		return NewListExtractors(apiClient.Extractors, out), nil
	case new(ClassifyCommand).Usage():
		return NewClassifyFilesCommandFromArgs(args, apiClient)
	}

	return nil, errors.New(fmt.Sprintf("Unknown command: %s", cmd))
}

var DefaultHttpClient = &http.Client{Timeout: time.Minute * 5}

func initSurf(args docopt.Opts) (*ch360.ApiClient, error) {

	appDir, err := config.NewAppDirectory()
	if err != nil {
		return nil, err
	}

	credentialsResolver := &CredentialsResolver{}

	clientId, clientSecret, err := credentialsResolver.ResolveFromArgs(args, appDir)

	if err != nil {
		return nil, err
	}

	return ch360.NewApiClient(DefaultHttpClient, ch360.ApiAddress, clientId, clientSecret), nil
}

func parseCommand(args docopt.Opts) (cmd string, err error) {
	verb, err := verbFromArgs(args)
	if err != nil {
		return
	}
	noun := nounFromArgs(args)

	if noun != "" {
		cmd = fmt.Sprintf("%s %s", verb, noun)
	} else {
		cmd = verb
	}

	return
}

func verbFromArgs(args docopt.Opts) (string, error) {
	supportedVerbs := []string{"login", "list", "create", "delete", "classify"}
	for _, verb := range supportedVerbs {
		if v, _ := args.Bool(verb); v {
			return verb, nil
		}
	}
	return "", errors.New("No supported verbs found.")
}

func nounFromArgs(args docopt.Opts) string {
	supportedNouns := []string{"classifier", "classifiers", "extractor", "extractors"}
	for _, noun := range supportedNouns {
		if v, _ := args.Bool(noun); v {
			return noun
		}
	}
	return ""
}
