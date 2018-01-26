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
  surf list extractors [options]
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
  -p, --progress                   : Show progress when classifying files. Only visible when
                                     redirecting stdout or in conjunction with -m or -o.
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

	exitOnErr(err)

	ctx, canceller := context.WithCancel(context.Background())

	go handleInterrupt(canceller)

	verb, err := verbFromArgs(args)
	exitOnErr(err)
	noun := nounFromArgs(args)

	var cmd string
	if noun != "" {
		cmd = fmt.Sprintf("%s %s", verb, noun)
	} else {
		cmd = verb
	}

	switch cmd {
	case "login":
		exitOnErr(doLogin(args))
	case "create classifier":
		exitOnErr(doCreateClassifier(args))
	case "delete classifier":
		exitOnErr(doDeleteClassifier(args))
	case "list classifiers":
		exitOnErr(doListClassifiers(args))
	case "list extractors":
		exitOnErr(doListExtractors(args))
	case "classify":
		exitOnErr(doClassifyFiles(ctx, args))
	}

}
func doListExtractors(args map[string]interface{}) error {
	clientId, clientSecret, err := resolveCredentials(args)

	if err != nil {
		return err
	}

	apiClient := ch360.NewApiClient(httpClient(), ch360.ApiAddress, clientId, clientSecret)

	_, err = commands.NewListExtractors(os.Stdout, apiClient.Extractors).Execute()

	return err
}

func doClassifyFiles(ctx context.Context, args map[string]interface{}) error {
	var (
		outputFormat       = argAsString(args, "--output-format")
		outputFilename     = argAsString(args, "--output-file")
		writeMultipleFiles = args["--multiple-files"].(bool)
		filePattern        = args["<file>"].(string)
		classifierName     = args["<classifier>"].(string)
		showProgress       = args["--progress"].(bool)
	)

	clientId, clientSecret, err := resolveCredentials(args)

	if err != nil {
		return err
	}

	apiClient := ch360.NewApiClient(httpClient(), ch360.ApiAddress, clientId, clientSecret)

	builder := resultsWriters.NewResultsWriterBuilder(outputFormat,
		writeMultipleFiles,
		outputFilename)

	resultsWriter, err := builder.Build()

	if err != nil {
		return err
	}

	// Only show progress if stdout is being redirected
	if !shouldShowProgressBar(writeMultipleFiles || outputFilename != "") {
		showProgress = false
	}

	progressHandler := progress.NewClassifyProgressHandler(resultsWriter, showProgress, os.Stderr)
	return commands.NewClassifyCommand(progressHandler,
		apiClient.Documents,
		apiClient.Documents,
		apiClient.Documents,
		apiClient.Documents,
		10).Execute(ctx, filePattern, classifierName)
}

func doListClassifiers(args map[string]interface{}) error {
	clientId, clientSecret, err := resolveCredentials(args)

	if err != nil {
		return err
	}

	apiClient := ch360.NewApiClient(httpClient(), ch360.ApiAddress, clientId, clientSecret)

	_, err = commands.NewListClassifiers(os.Stdout, apiClient.Classifiers).Execute()

	return err
}

func doDeleteClassifier(args map[string]interface{}) error {
	classifierName := args["<name>"].(string)

	clientId, clientSecret, err := resolveCredentials(args)

	if err != nil {
		return err
	}

	apiClient := ch360.NewApiClient(httpClient(), ch360.ApiAddress, clientId, clientSecret)

	return commands.NewDeleteClassifier(os.Stdout, apiClient.Classifiers).Execute(classifierName)
}

func resolveCredentials(args map[string]interface{}) (clientId, clientSecret string, err error) {
	var (
		clientIdArg     = argAsString(args, "--client-id")
		clientSecretArg = argAsString(args, "--client-secret")
	)

	appDir, err := appDirectory()

	if err != nil {
		return
	}

	// Get credentials from configuration
	resolver := &commands.CredentialsResolver{}
	clientId, clientSecret, err = resolver.Resolve(clientIdArg, clientSecretArg, appDir)

	return
}

func doCreateClassifier(args map[string]interface{}) error {
	clientId, clientSecret, err := resolveCredentials(args)

	if err != nil {
		return err
	}

	apiClient := ch360.NewApiClient(httpClient(), ch360.ApiAddress, clientId, clientSecret)

	classifierName := args["<name>"].(string)
	samplesPath := args["<samples-zip>"].(string)

	return commands.NewCreateClassifier(
		os.Stdout,
		apiClient.Classifiers,
		apiClient.Classifiers,
		apiClient.Classifiers,
	).Execute(classifierName, samplesPath)

}

func appDirectory() (*config.AppDirectory, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	appDirectory := config.NewAppDirectory(homedir)
	return appDirectory, nil
}

func doLogin(args map[string]interface{}) error {
	var (
		clientId     = argAsString(args, "--client-id")
		clientSecret = argAsString(args, "--client-secret")
		err          error
	)

	if clientId == "" {
		fmt.Print("Client Id: ")
		clientId, err = readSecretFromConsole()
		if err != nil {
			return err
		}
	}

	if clientSecret == "" {
		fmt.Print("Client Secret: ")
		clientSecret, err = readSecretFromConsole()
		if err != nil {
			return err
		}
	}

	responseChecker := &response.ErrorChecker{}

	tokenRetriever := auth.NewHttpTokenRetriever(clientId,
		clientSecret,
		httpClient(),
		ch360.ApiAddress,
		responseChecker)

	appDir, err := appDirectory()

	if err != nil {
		return err
	}

	return commands.NewLogin(os.Stdout,
		appDir,
		tokenRetriever).Execute(clientId, clientSecret)

}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Minute * 5,
	}
}

func shouldShowProgressBar(writingToFile bool) bool {
	return writingToFile || isRedirected(os.Stdout.Fd())
}

func isRedirected(fd uintptr) bool {
	return !isatty.IsTerminal(fd) &&
		!isatty.IsCygwinTerminal(fd)
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

func verbFromArgs(args map[string]interface{}) (string, error) {
	supportedVerbs := []string{"login", "list", "create", "delete", "classify"}
	for _, verb := range supportedVerbs {
		if args[verb].(bool) {
			return verb, nil
		}
	}
	return "", errors.New("No supported verbs found.")
}

func nounFromArgs(args map[string]interface{}) string {
	supportedNouns := []string{"classifier", "classifiers", "extractors"}
	for _, noun := range supportedNouns {
		if args[noun].(bool) {
			return noun
		}
	}
	return ""
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
