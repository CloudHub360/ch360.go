package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/ioutils"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	//	usage := `surf - the official command line client for waives.io.
	//
	//Usage:
	//  surf login [options]
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
		globalFlags = config.GlobalFlags{}

		app = kingpin.New("surf", "surf - the official command line client for waives.io.").
			Version(ch360.Version)

		list            = app.Command("list", "List waives resources.")
		listModules     = list.Command("modules", "List all available extractor modules.")
		listClassifiers = list.Command("classifiers", "List all available classifiers.")
		listExtractors  = list.Command("extractors", "List all available extractors.")

		upload = app.Command("upload", "Upload waives resources.")

		uploadClassifier     = upload.Command("classifier", "Upload waives classifier (.clf file).")
		uploadClassifierName = uploadClassifier.Arg("name", "The name of the new classifier.").Required().String()
		uploadClassifierFile = uploadClassifier.Arg("classifier-file", "The trained classifier file.").Required().File()

		deleteCmd           = app.Command("delete", "Delete waives resources.")
		deleteExtractor     = deleteCmd.Command("extractor", "Delete waives extractor.")
		deleteExtractorName = deleteExtractor.Arg("name", "The name of the extractor to delete.").Required().String()

		deleteClassifier     = deleteCmd.Command("classifier", "Delete waives classifier.")
		deleteClassifierName = deleteClassifier.Arg("name", "The name of the classifier to delete.").Required().String()

		create               = app.Command("create", "Create waives resources.")
		createClassifier     = create.Command("classifier", "Create waives classifier from a set of samples.")
		createClassifierName = createClassifier.Arg("name", "The name of the new classifier.").Required().String()
		createClassifierFile = createClassifier.Arg("samples-zip", "The zip file containing training samples.").Required().File()

		createExtractor = create.Command("extractor", "Create waives extractor.")

		createExtractorFromModules = createExtractor.Command("from-modules", "Create waives extractor from a set of modules.")

		createExtractorFromModulesName = createExtractorFromModules.Arg("name", "The name of the new extractor.").Required().String()
		createExtractorFromModulesIds  = createExtractorFromModules.Arg("module-ids",
			"The module ids to create the extractor from.").Required().Strings()

		createExtractorFromTemplate     = createExtractor.Command("from-template", "The extractor template to create the extractor from.")
		createExtractorFromTemplateName = createExtractorFromTemplate.Arg("name", "The name of the new extractor.").
						Required().String()
		createExtractorFromTemplateFile = createExtractorFromTemplate.Arg("template-file", "The extraction template file (json).").
						Required().File()

		createExtractorTemplate        = create.Command("extractor-template", "Create an extractor template from the provided module ids")
		createExtractorTemplateModules = createExtractorTemplate.Arg("module-ids", "The module IDs to include in the template").
						Required().Strings()
	)
	ctx, canceller := context.WithCancel(context.Background())
	go handleInterrupt(canceller)

	commands.ConfigureLoginCommand(ctx, app, &globalFlags)
	commands.ConfigureUploadExtractorCommand(ctx, upload, &globalFlags)
	commands.ConfigureReadCommand(ctx, app, &globalFlags)
	commands.ConfigureExtractCommand(ctx, app, &globalFlags)
	commands.ConfigureClassifyCommand(ctx, app, &globalFlags)

	app.Flag("client-id", "Client ID").Short('i').
		Short('i').
		StringVar(&globalFlags.ClientId)
	app.Flag("client-secret", "Client secret").Short('s').
		Short('s').
		StringVar(&globalFlags.ClientSecret)
	app.Flag("log-http", "Log HTTP requests and responses as they happen, "+
		"to a file.").
		OpenFileVar(&globalFlags.LogHttp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	app.Flag("multiple-files",
		"Write results output to multiple files with the same basename as the input").
		Short('m').
		BoolVar(&globalFlags.MultiFileOut)
	app.Flag("output-file", "Write all results to the specified file").
		Short('o').
		StringVar(&globalFlags.OutputFile)

	defer ioutils.TryClose(globalFlags.LogHttp)

	parsedCommand := kingpin.MustParse(app.Parse(os.Args[1:]))

	var cmd commands.Command

	apiClient, err := initApiClient(globalFlags.ClientId, globalFlags.ClientSecret, globalFlags.LogHttp)
	exitOnErr(err)

	switch parsedCommand {
	case listModules.FullCommand():
		cmd = commands.NewListModules(apiClient.Modules, os.Stdout)
	case listClassifiers.FullCommand():
		cmd = commands.NewListClassifiers(apiClient.Classifiers, os.Stdout)
	case listExtractors.FullCommand():
		cmd = commands.NewListExtractors(apiClient.Extractors, os.Stdout)
	case uploadClassifier.FullCommand():
		cmd = commands.NewUploadClassifier(os.Stdout, apiClient.Classifiers, *uploadClassifierName, *uploadClassifierFile)
		defer (*uploadClassifierFile).Close()
	case deleteExtractor.FullCommand():
		cmd = commands.NewDeleteExtractor(*deleteExtractorName, os.Stdout, apiClient.Extractors)
	case deleteClassifier.FullCommand():
		cmd = commands.NewDeleteClassifier(*deleteClassifierName, os.Stdout, apiClient.Classifiers)
	case createClassifier.FullCommand():
		cmd = commands.NewCreateClassifier(os.Stdout, apiClient.Classifiers,
			apiClient.Classifiers, apiClient.Classifiers, *createClassifierName, *createClassifierFile)
		defer (*createClassifierFile).Close()
	case createExtractorFromModules.FullCommand():
		cmd = commands.NewCreateExtractorFromModules(os.Stdout, apiClient.Extractors,
			*createExtractorFromModulesName, *createExtractorFromModulesIds)
	case createExtractorFromTemplate.FullCommand():
		cmd, err = commands.NewCreateExtractorFromTemplate(os.Stdout, apiClient.Extractors,
			*createExtractorFromTemplateName, *createExtractorFromTemplateFile)

		exitOnErr(err)

	case createExtractorTemplate.FullCommand():
		out := os.Stdout

		cmd = commands.NewCreateExtractorTemplate(*createExtractorTemplateModules,
			apiClient.Modules, out)

	}

	if cmd != nil {
		exitOnErr(cmd.Execute(ctx))
	}
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
