package commands

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/services"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/progress"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type ClassifyArgs struct {
	classifierName string
	outputFormat   string
	filePatterns   []string
}

func ConfigureClassifyCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	classifyArgs := &ClassifyArgs{}

	cmd := app.
		Command("classify", "Perform content classification on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			// execute the command
			return ExecuteClassify(ctx, classifyArgs, globalFlags)
		})

	cmd.Flag("format", "The output format. Allowed values: table, csv, json [default: table].").
		Short('f').
		Default("table").
		EnumVar(&classifyArgs.outputFormat, "table", "csv", "json")

	cmd.Arg("classifier-name", "The name of the classifier to use.").
		Required().
		StringVar(&classifyArgs.classifierName)

	cmd.Arg("files", "The files to read.").
		Required().
		StringsVar(&classifyArgs.filePatterns)
}

// ExecuteClassify is the main entry point for the 'classify' command. It has to do a lot of
// setup / instantiation before actually performing classification on the specified files.
func ExecuteClassify(ctx context.Context, classifyArgs *ClassifyArgs,
	globalFlags *config.GlobalFlags) error {

	resultsWriter, err := resultsWriters.NewClassificationResultsWriter(globalFlags.MultiFileOut,
		globalFlags.OutputFile,
		classifyArgs.outputFormat)

	if err != nil {
		return err
	}

	progressHandler := progress.NewProgressHandler(resultsWriter,
		globalFlags.ShowProgress, os.Stderr)

	filePaths, err := GlobMany(classifyArgs.filePatterns)

	if err != nil {
		return err
	}

	client, err := initApiClient(globalFlags.ClientId,
		globalFlags.ClientSecret,
		globalFlags.LogHttp)

	if err != nil {
		return err
	}

	fileClassifier := ch360.NewFileClassifier(client.Documents, client.Documents, client.Documents)

	return services.NewParallelClassificationService(fileClassifier, client.Documents,
		progressHandler).ClassifyAll(ctx, filePaths, classifyArgs.classifierName)
}
