package main

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/progress"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"github.com/CloudHub360/ch360.go/response"
	"github.com/docopt/docopt-go"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	usage := `surf - the official command line client for waives.io.

Usage:
  surf login [options]
  surf create classifier <name> <samples-zip> [options]
  surf delete classifier <name> [options]
  surf list classifiers [options]
  surf classify <file> <classifier> [options]
  surf -h | --help
  surf -v | --version

Options:
  -h, --help                       : Show this help message
  -v, --version                    : Show application version
  -i, --client-id <id>             : Client ID
  -s, --client-secret <secret>     : Client secret
  -f, --output-format <format>     : Output format for classification results. Allowed values: table, csv, json [default: table]
  -o, --output-file <file>         : Write all results to the specified file
  -m, --multiple-files             : Write results output to multiple files with the same
                                   : basename as the input
  -p, --progress                   : Show progress
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
			if clientId, err = readSecretFromConsole(); err != nil {
				os.Exit(1)
			}
		}

		if clientSecret == "" {
			fmt.Print("Client Secret: ")
			if clientSecret, err = readSecretFromConsole(); err != nil {
				os.Exit(1)
			}
		}

		responseChecker := &response.ErrorChecker{}
		tokenRetriever := auth.NewHttpTokenRetriever(clientId, clientSecret, httpClient, ch360.ApiAddress, responseChecker)
		err = commands.NewLogin(os.Stdout, appDirectory, tokenRetriever).Execute(clientId, clientSecret)
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
			commands.NewDeleteClassifier(os.Stdout, apiClient.Classifiers),
		).Execute(classifierName, samplesPath)
		if err != nil {
			os.Exit(1)
		}
	} else if args["delete"].(bool) {
		classifierName := args["<name>"].(string)

		fmt.Printf("Deleting classifier '%s'... ", classifierName)
		err = commands.NewDeleteClassifier(os.Stdout, apiClient.Classifiers).Execute(classifierName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("[OK]")
	} else if args["list"].(bool) {
		var classifiers ch360.ClassifierList
		classifiers, err = commands.NewListClassifiers(os.Stdout, apiClient.Classifiers).Execute()

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

		var (
			outputFormat       = argAsString(args, "--output-format")
			outputFilename     = argAsString(args, "--output-file")
			writeMultipleFiles = args["--multiple-files"].(bool)
			filePattern        = args["<file>"].(string)
			classifierName     = args["<classifier>"].(string)
			showProgress       = args["--progress"].(bool)
		)

		builder := resultsWriters.NewResultsWriterBuilder(outputFormat, writeMultipleFiles, outputFilename)
		resultsWriter, err := builder.Build()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// Only show progress if stdout is being redirected
		if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			showProgress = false
		}

		progressHandler := progress.NewClassifyProgressHandler(resultsWriter, showProgress, os.Stderr)
		err = commands.NewClassifyCommand(progressHandler,
			apiClient.Documents,
			apiClient.Documents,
			apiClient.Documents,
			apiClient.Documents,
			10).Execute(ctx, filePattern, classifierName)
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
		if err != commands.ConsoleSecretReaderErrCancelled {
			fmt.Println(err.Error())
		}
		return "", err
	}
	return secret, nil
}
