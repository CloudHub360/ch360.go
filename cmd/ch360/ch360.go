package main

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/response"
	"github.com/docopt/docopt-go"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	usage := `CloudHub360 command-line tool.

Usage:
  ch360 login [options]
  ch360 create classifier <name> <samples-zip> [options]
  ch360 delete classifier <name> [options]
  ch360 list classifiers [options]
  ch360 classify <file> <classifier> [options]
  ch360 -h | --help
  ch360 -v | --version

Options:
  -h, --help                       : Show this help message
  -v, --version                    : Show application version
  -i, --client-id <id>             : Client ID
  -s, --client-secret <secret>     : Client secret
  -f, --output-format <format>     : Output format for classification results. Allowed values: table, csv [default: table]
`

	filenameExamples := `
Filename and glob pattern examples:
  file1.pdf        : Specific file
  *.*              : All files in the current folder
  *.pdf            : All PDFs in the current folder
  foo/*.tif        : All TIFs in folder foo
  bar/**/*.*       : All files in subfolders of folder bar`

	// Replace slashes with OS-specific path separators
	usage = usage + filepath.FromSlash(filenameExamples)

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

	homedir, err := homedir.Dir()
	if err != nil {
		fmt.Println(errors.New(fmt.Sprintf("Could not determine home directory. Details: %s", err.Error())))
		os.Exit(1)
	}

	appDirectory := config.NewAppDirectory(homedir)

	ctx, canceller := context.WithCancel(context.Background())

	go handleInterrupt(canceller)

	if args["login"].(bool) {
		if clientId == "" {
			fmt.Print("Client Id: ")
			clientId, err = readSecretFromConsole()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		if clientSecret == "" {
			fmt.Print("Client Secret: ")
			clientSecret, err = readSecretFromConsole()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
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
		samplesPath := args["<samples-zip>"].(string)
		fmt.Printf("Creating classifier '%s'... ", classifierName)
		err = commands.NewCreateClassifier(
			os.Stdout,
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
	} else if args["classify"].(bool) {
		filePattern := args["<file>"].(string)
		classifierName := args["<classifier>"].(string)

		outputFormat := argAsString(args, "--output-format")
		var writer commands.ClassifyResultsWriter
		switch outputFormat {
		case "table":
			writer = commands.NewTableClassifyResultsWriter(os.Stdout)
		case "csv":
			writer = commands.NewCSVClassifyResultsWriter(os.Stdout)
		default:
			// DocOpt doesn't do validation of these values for us, so we need to catch invalid values here
			fmt.Println(fmt.Sprintf("Unknown output format '%s'. Allowed values are: csv, table.", outputFormat))
			os.Exit(1)
		}

		err = commands.NewClassifyCommand(writer, os.Stdout, apiClient.Documents, 10).Execute(ctx, filePattern, classifierName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func handleInterrupt(canceller context.CancelFunc) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan // ctrl-c received
	canceller()
}

func argAsString(args map[string]interface{}, name string) string {
	var result string = ""
	if args[name] != nil {
		result = args[name].(string)
	}

	return result
}

func readSecretFromConsole() (string, error) {
	secret, err := (&commands.ConsoleSecretReader{}).Read()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return secret, nil
}
