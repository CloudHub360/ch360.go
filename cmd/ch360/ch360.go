package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"time"
	"github.com/CloudHub360/ch360.go/ch360"
)

func main() {
	usage := `CloudHub360 command-line tool.

Usage:
  ch360 create classifier <name> --id=<id> --secret=<secret>
  ch360 -h | --help
  ch360 --version

Options:
  -h --help          Show this help message.
  --version          Show version.
  --id=<id>          API Client ID
  --secret=<secret>  API Client secret`

	args, err := docopt.Parse(usage, nil, true, ch360.Version, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	id := args["--id"].(string)
	secret := args["--secret"].(string)
	classifierName := args["<name>"].(string)

	var httpClient = &http.Client{
		Timeout: time.Minute * 5,
	}

	apiClient := ch360.NewApiClient(httpClient, ch360.ApiAddress, id, secret)
	err = apiClient.Classifiers.CreateClassifier(classifierName)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Created classifier '%s'.\n", classifierName)
}