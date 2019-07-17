package main

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/ioutils"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//	usage := `surf - the official command line client for waives.io.
	//
	//Usage:
	//  surf login [options]
	//  surf ` + new(commands.ListClassifiers).Usage() + ` [options]
	//  surf ` + new(commands.UploadClassifier).Usage() + ` <name> <classifier-file> [options]
	//  surf ` + new(commands.CreateClassifier).Usage() + ` <name> <samples-zip> [options]
	//  surf ` + new(commands.DeleteClassifier).Usage() + ` <name> [options]
	//  surf ` + new(commands.ListExtractors).Usage() + ` [options]
	//  surf ` + new(commands.UploadExtractor).Usage() + ` <name> <config-file> [options]
	//  surf ` + new(commands.CreateExtractor).Usage() + ` <name> --from-template=<template> [options]
	//  surf ` + new(commands.CreateExtractor).Usage() + ` <name> <module-ids>... [options]
	//  surf ` + new(commands.CreateExtractorTemplate).Usage() + ` <module-ids>... [options]
	//  surf ` + new(commands.DeleteExtractor).Usage() + ` <name> [options]
	//  surf ` + new(commands.ListModules).Usage() + ` [options]
	//  surf ` + new(commands.ClassifyCommand).Usage() + ` <file> <classifier> [options]
	//  surf ` + new(commands.Extract).Usage() + ` <file> <extractor> [options]
	//  surf ` + new(commands.Read).Usage() + ` <file> (pdf|txt|wvdoc) [options]
	//  surf -h | --help
	//  surf -v | --version
	//
	//Options:
	//  -h, --help                       : Show this help message
	//  -v, --version                    : Show application version
	//  -i, --client-id <id>             : Client ID
	//  -s, --client-secret <secret>     : Client secret
	//  -f, --output-format <format>     : Output format for classification and extraction results.
	//                                     Allowed values: table, csv, json [default: table]
	//  -o, --output-file <file>         : Write all results to the specified file
	//  -m, --multiple-files             : Write results output to multiple files with the same
	//                                   : basename as the input
	//  -p, --progress                   : Show progress when classifying files. Only visible when
	//                                     redirecting stdout or in conjunction with -m or -o.
	//  -t, --from-template <template>   : The extractor modules template to use when creating an
	//                                     extractor from modules.
	//  --log-http <file>                : Log HTTP requests and responses as they happen, to a file.
	//`
	//
	//	filenameExamples := `
	//Filename and glob pattern examples:
	//  file1.pdf        : Specific file
	//  *.*              : All files in the current folder
	//  *.pdf            : All PDFs in the current folder
	//  foo/*.tif        : All TIFs in folder foo
	//  bar/**/*.*       : All files in subfolders of folder bar`
	//
	//	// Replace slashes with OS-specific path separators
	//	usage = usage + filepath.FromSlash(filenameExamples)
	//
	//	args, err := docopt.ParseArgs(usage, nil, ch360.Version)
	//	exitOnErr(err)
	//
	//	runParams, err := config.NewRunParamsFromArgs(args)
	//	exitOnErr(err)
	//
	//	ctx, canceller := context.WithCancel(context.Background())
	//	go handleInterrupt(canceller)
	//
	//	appDir, err := config.NewAppDirectory()
	//	exitOnErr(err)
	//
	//	var (
	//		cmd       commands.Command
	//		apiClient *ch360.ApiClient
	//	)
	//
	//	if login, _ := args.Bool("login"); login {
	//		tokenRetriever := ch360.NewTokenRetriever(DefaultHttpClient, ch360.ApiAddress)
	//		cmd = commands.NewLoginFrom(runParams, os.Stdout, appDir, tokenRetriever)
	//	} else {
	//		apiClient, err = initApiClient(runParams)
	//		exitOnErr(err)
	//		cmd, err = commands.CommandFor(runParams, apiClient)
	//	}
	//	exitOnErr(err)
	//
	//	exitOnErr(cmd.Execute(ctx))

	var (
		app          = kingpin.New("surf", "surf - the official command line client for waives.io.")
		clientId     = app.Flag("client-id", "Client ID").Short('i').String()
		clientSecret = app.Flag("client-secret", "Client secret").Short('s').String()
		logHttp      = app.Flag("log-http", "Log HTTP requests and responses as they happen, to a file.").File()

		list            = app.Command("list", "List waives resources.")
		listModules     = list.Command("modules", "List all available extractor modules.")
		listClassifiers = list.Command("classifiers", "List all available classifiers.")
		listExtractors  = list.Command("extractors", "List all available extractors.")

		upload = app.Command("upload", "Upload waives resources.")

		uploadExtractor     = upload.Command("extractor", "Upload waives extractor (.fpxlc file).")
		uploadExtractorName = uploadExtractor.Arg("name", "The name of the new extractor.").String()
		uploadExtractorFile = uploadExtractor.Arg("config-file", "The extraction configuration file.").File()

		uploadClassifier     = upload.Command("classifier", "Upload waives classifier (.clf file).")
		uploadClassifierName = uploadClassifier.Arg("name", "The name of the new classifier.").String()
		uploadClassifierFile = uploadClassifier.Arg("classifier-file", "The trained classifier file.").File()

		login = app.Command("login", "Connect surf to your account.")
	)
	defer ioutils.TryClose(*logHttp)

	parsedCommand := kingpin.MustParse(app.Parse(os.Args[1:]))

	appDir, err := config.NewAppDirectory()
	exitOnErr(err)

	ctx, canceller := context.WithCancel(context.Background())
	go handleInterrupt(canceller)

	var cmd commands.Command

	if parsedCommand == login.FullCommand() {
		// special case for login, it doesn't need the api client to be created
		tokenRetriever := ch360.NewTokenRetriever(DefaultHttpClient, ch360.ApiAddress)
		cmd = commands.NewLogin(os.Stdout, appDir, tokenRetriever, *clientId, *clientSecret)
	} else {

		apiClient, err := initApiClient(*clientId, *clientSecret, *logHttp)
		exitOnErr(err)

		switch parsedCommand {
		case listModules.FullCommand():
			cmd = commands.NewListModules(apiClient.Modules, os.Stdout)
		case listClassifiers.FullCommand():
			cmd = commands.NewListClassifiers(apiClient.Classifiers, os.Stdout)
		case listExtractors.FullCommand():
			cmd = commands.NewListExtractors(apiClient.Extractors, os.Stdout)
		case uploadExtractor.FullCommand():
			cmd = commands.NewUploadExtractor(os.Stdout, apiClient.Extractors, *uploadExtractorName, *uploadExtractorFile)
			defer (*uploadExtractorFile).Close()
		case uploadClassifier.FullCommand():
			cmd = commands.NewUploadClassifier(os.Stdout, apiClient.Classifiers, *uploadClassifierName, *uploadClassifierFile)
			defer (*uploadClassifierFile).Close()
		}
	}

	exitOnErr(cmd.Execute(ctx))
}

func handleInterrupt(canceller context.CancelFunc) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan // ctrl-c received
	canceller()
}

func exitOnErr(err error) {
	if err != nil && err != context.Canceled {

		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func initApiClient(clientIdFlag, clientSecretFlag string, logHttpFile *os.File) (*ch360.ApiClient, error) {
	appDir, err := config.NewAppDirectory()
	if err != nil {
		return nil, err
	}

	credentialsResolver := &commands.CredentialsResolver{}

	clientId, clientSecret, err := credentialsResolver.Resolve(clientIdFlag, clientSecretFlag, appDir)

	if err != nil {
		return nil, err
	}

	var logSink io.Writer = nil
	if logHttpFile != nil {
		logSink = logHttpFile
	}
	return ch360.NewApiClient(DefaultHttpClient, ch360.ApiAddress, clientId, clientSecret, logSink), nil
}

var DefaultHttpClient = &http.Client{Timeout: time.Minute * 2}
