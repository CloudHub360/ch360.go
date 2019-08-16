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

type ExtractArgs struct {
	extractorName string
	outputFormat  string
	filePatterns  []string
}

func ConfigureExtractCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	extractArgs := &ExtractArgs{}

	cmd := app.
		Command("extract", "Perform data extraction on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			// execute the command
			return ExecuteExtract(ctx, extractArgs, globalFlags)
		})

	cmd.Flag("format", "The output format. Allowed values: table, csv, json [default: table].").
		Short('f').
		Default("table").
		EnumVar(&extractArgs.outputFormat, "table", "csv", "json")

	cmd.Arg("extractor-name", "The name of the extractor to use.").
		Required().
		StringVar(&extractArgs.extractorName)

	cmd.Arg("files", "The files to read.").
		Required().
		StringsVar(&extractArgs.filePatterns)
}

// ExecuteExtract is the main entry point for the 'extract' command. It has to do a lot of
// setup / instantiation before actually performing data extraction on the specified files.
func ExecuteExtract(ctx context.Context, extractArgs *ExtractArgs,
	globalFlags *config.GlobalFlags) error {

	resultsWriter, err := resultsWriters.NewExtractionResultsWriter(globalFlags.MultiFileOut,
		globalFlags.OutputFile,
		extractArgs.outputFormat)

	if err != nil {
		return err
	}

	progressHandler := progress.NewProgressHandler(resultsWriter,
		globalFlags.ShowProgress, os.Stderr)

	filePaths, err := GlobMany(extractArgs.filePatterns)

	if err != nil {
		return err
	}

	client, err := initApiClient(globalFlags.ClientId,
		globalFlags.ClientSecret,
		globalFlags.LogHttp)

	if err != nil {
		return err
	}

	fileExtractor := ch360.NewFileExtractor(client.Documents, client.Documents, client.Documents)

	return services.NewParallelExtractorService(fileExtractor, client.Documents, progressHandler).
		ExtractAll(ctx, filePaths, extractArgs.extractorName)
}
