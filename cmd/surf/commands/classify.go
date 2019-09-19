package commands

import (
	"context"
	"github.com/pkg/errors"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/cmd/surf/services"
	"github.com/waives/surf/config"
	"github.com/waives/surf/output/progress"
	"github.com/waives/surf/output/resultsWriters"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type classifyArgs struct {
	classifierName string
	outputFormat   string
	filePatterns   []string
}

//go:generate mockery -name "ClassificationService"
type ClassificationService interface {
	ClassifyAll(ctx context.Context, files []string, classifierName string) error
}

type ClassifyCmd struct {
	ClassificationService ClassificationService
	FilePaths             []string
	ClassifierName        string
}

func ConfigureClassifyCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	classifyArgs := &classifyArgs{}
	classifyCmd := &ClassifyCmd{}

	classifyCli := app.
		Command("classify", "Perform content classification on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := classifyCmd.initWithArgs(classifyArgs, globalFlags)
			if err != nil {
				return err
			}

			return classifyCmd.Execute(ctx)
		})

	classifyCli.Flag("format", "The output format. Allowed values: table, csv, json [default: table].").
		Short('f').
		Default("table").
		EnumVar(&classifyArgs.outputFormat, "table", "csv", "json")

	classifyCli.Arg("classifier-name", "The name of the classifier to use.").
		Required().
		StringVar(&classifyArgs.classifierName)

	classifyCli.Arg("files", "The files to read.").
		Required().
		StringsVar(&classifyArgs.filePatterns)

	addFileHandlingFlagsTo(globalFlags, classifyCli)
}

// ExecuteClassify is the main entry point for the 'classify' command.
func (cmd *ClassifyCmd) Execute(ctx context.Context) error {
	err := cmd.ClassificationService.ClassifyAll(ctx, cmd.FilePaths, cmd.ClassifierName)
	return errors.Wrap(err, "classification failed")
}

func (cmd *ClassifyCmd) initWithArgs(args *classifyArgs, flags *config.GlobalFlags) error {
	resultsWriter, err := resultsWriters.NewClassificationResultsWriter(flags.MultiFileOut,
		flags.OutputFile,
		args.outputFormat)

	if err != nil {
		return err
	}

	progressHandler := progress.NewProgressHandler(resultsWriter,
		flags.ShowProgress, os.Stderr)

	cmd.FilePaths, err = GlobMany(args.filePatterns)

	if err != nil {
		return err
	}

	client, err := initApiClient(flags.ClientId,
		flags.ClientSecret,
		flags.LogHttp)

	if err != nil {
		return err
	}

	fileClassifier := ch360.NewFileClassifier(client.Documents, client.Documents, client.Documents)

	cmd.ClassificationService = services.NewParallelClassificationService(fileClassifier,
		client.Documents,
		progressHandler)
	cmd.ClassifierName = args.classifierName

	return nil
}
