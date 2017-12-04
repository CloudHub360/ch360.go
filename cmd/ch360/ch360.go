package main

import (
	"fmt"
	"github.com/CloudHub360/ch360.go"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"time"
	ch3602 "github.com/CloudHub360/ch360.go/ch360"
)

func main() {
	usage := `CloudHub360 command-line tool.

Usage:
  ch360 -h | --help
  ch360 --version
  ch360 --id=<id> --secret=<secret>

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

	var httpClient = &http.Client{
		Timeout: time.Minute * 5,
	}

	tokenGetter := auth.NewHttpTokenRetriever(id, secret, httpClient, ch3602.ApiAddress)

	apiClient := ch3602.NewApiClient(httpClient, ch3602.ApiAddress, tokenGetter)
	err = apiClient.Classifiers.CreateClassifier("myclassifier")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("Created classifier 'myclassifier'.")
}
