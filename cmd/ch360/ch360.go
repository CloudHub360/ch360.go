package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"time"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/auth"
	buildvars "github.com/CloudHub360/ch360.go"
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

	args, err := docopt.Parse(usage, nil, true, buildvars.Version, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	id := args["--id"].(string)
	secret := args["--secret"].(string)

	var httpClient = &http.Client{
		Timeout: time.Minute * 5,
	}

	tokenGetter := auth.NewHttpTokenRetriever(id, secret, httpClient, ch360.ApiAddress)

	apiClient := ch360.NewApiClient(httpClient, ch360.ApiAddress, tokenGetter)
	err = apiClient.Classifiers.CreateClassifier("myclassifier")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("Created classifier 'myclassifier'.")
}