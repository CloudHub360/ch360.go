package main

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"os/user"
	"time"
)

func main() {
	usage := `CloudHub360 command-line tool.

Usage:
  ch360 login --client-id=<id> [--client-secret=<secret>]
  ch360 create classifier <name> --client-id=<id> --client-secret=<secret> --samples-zip=<path>
  ch360 delete classifier <name> --client-id=<id> --client-secret=<secret>
  ch360 list classifiers --client-id=<id> --client-secret=<secret>
  ch360 -h | --help
  ch360 --version

Options:
  -h --help          Show this help message.
  --version          Show version.
  --client-id=<id>          API Client ID
  --client-secret=<secret>  API Client secret`

	args, err := docopt.Parse(usage, nil, true, ch360.Version, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	if args["login"].(bool) {
		id := args["--client-id"].(string)

		var secret = ""
		if args["--client-secret"] != nil {
			secret = args["--client-secret"].(string)
		}

		user, err := user.Current()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		appDirectory := config.NewAppDirectory(user.HomeDir)
		err = commands.NewLogin(appDirectory, &commands.ConsoleSecretReader{}).Execute(id, secret)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	var httpClient = &http.Client{
		Timeout: time.Minute * 5,
	}

	id := args["--client-id"].(string)
	secret := args["--client-secret"].(string)
	apiClient := ch360.NewApiClient(httpClient, ch360.ApiAddress, id, secret)

	if args["create"].(bool) {
		classifierName := args["<name>"].(string)
		samplesPath := args["--samples-zip"].(string)
		fmt.Printf("Creating classifier '%s'... ", classifierName)
		err = commands.NewCreateClassifier(
			apiClient.Classifiers,
			commands.NewDeleteClassifier(apiClient.Classifiers),
		).Execute(classifierName, samplesPath)
		if err != nil {
			os.Exit(1)
		}
	} else if args["delete"].(bool) {
		classifierName := args["<name>"].(string)

		fmt.Printf("Deleting classifier '%s'... ", classifierName)
		err = commands.NewDeleteClassifier(apiClient.Classifiers).Execute(classifierName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println("[OK]")
	} else if args["list"].(bool) {
		var classifiers ch360.ClassifierList
		classifiers, err = commands.NewListClassifiers(apiClient.Classifiers).Execute()

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		if !classifiers.Any() {
			fmt.Println("No classifiers found.")
		}

		for _, classifier := range classifiers {
			fmt.Println(classifier.Name)
		}
	}
}
