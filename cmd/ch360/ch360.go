package main

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/response"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"os/user"
	"time"
)

func main() {
	usage := `CloudHub360 command-line tool.

Usage:
  ch360 login [options]
  ch360 create classifier <name> --samples-zip=<path> [options]
  ch360 delete classifier <name> [options]
  ch360 list classifiers [options]
  ch360 -h | --help
  ch360 --version

Options:
  -h --help          Show this help message.
  --version          Show version.
  --client-id=<id>          API Client ID
  --client-secret=<secret>  API Client secret`

	args, err := docopt.Parse(usage, nil, true, ch360.Version, false)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}

	clientId := argAsString(args, "--client-id")
	clientSecret := argAsString(args, "--client-secret")
	if clientSecret == "" {
		clientSecret, err = readSecret()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	user, err := user.Current()
	if err != nil {
		fmt.Println(os.Stderr, err.Error())
		os.Exit(1)
	}
	appDirectory := config.NewAppDirectory(user.HomeDir)

	if args["login"].(bool) {
		if clientId == "" {
			//TODO: Prompt for client id if not specified
			fmt.Println("Please specify your API Client Id with the --client-id parameter")
			os.Exit(1)
		}

		responseChecker := &response.ErrorChecker{}
		tokenRetriever := auth.NewHttpTokenRetriever(clientId, clientSecret, httpClient, ch360.ApiAddress, responseChecker)
		err = commands.NewLogin(appDirectory, tokenRetriever).Execute(clientId, clientSecret)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	// Get credentials from configuration
	resolver := &commands.CredentialsResolver{}
	clientId, clientSecret, err = resolver.Resolve(clientId, clientSecret, appDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	apiClient := ch360.NewApiClient(httpClient, ch360.ApiAddress, clientId, clientSecret)

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
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("[OK]")
	} else if args["list"].(bool) {
		var classifiers ch360.ClassifierList
		classifiers, err = commands.NewListClassifiers(apiClient.Classifiers).Execute()

		if err != nil {
			fmt.Println(err.Error())
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

func argAsString(args map[string]interface{}, name string) string {
	var result string = ""
	if args[name] != nil {
		result = args[name].(string)
	}

	return result
}

func readSecret() (string, error) {
	fmt.Print("API Client Secret: ")
	secret, err := (&commands.ConsoleSecretReader{}).Read()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return secret, nil
}
