package main

import (
	//"fmt"

	"github.com/CloudHub360/ch360.go"
	"github.com/docopt/docopt-go"
	"fmt"
	"github.com/CloudHub360/ch360.go/authtoken"
	"net/http"
	"time"
	"os"
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

	tokenGetter := authtoken.NewHttpGetter(id, secret, httpClient, ch360.ApiAddress)

	apiClient := ch360.NewApiClient(httpClient, ch360.ApiAddress, tokenGetter)
	err = apiClient.CreateClassifier("myclassifier")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("Created classifier 'myclassifier'.")
}
