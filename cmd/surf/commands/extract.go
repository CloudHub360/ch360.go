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

type extractArgs struct {
	extractorName string
	outputFormat  string
	filePatterns  []string
}

type ExtractCmd struct {
	ExtractorName     string
	FilePaths         []string
	ExtractionService ExtractionService
}

//go:generate mockery -name "ExtractionService"
type ExtractionService interface {
	ExtractAll(ctx context.Context, files []string,
		extractorName string) error
}

func ConfigureExtractCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	args := &extractArgs{}
	cmd := &ExtractCmd{}
	extractCli := app.
		Command("extract", "Perform data extraction on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := cmd.initWithArgs(args, globalFlags)
			if err != nil {
				return err
			}
			return cmd.Execute(ctx)
		})

	extractCli.Flag("format", "The output format. Allowed values: table, csv, json [default: table].").
		Short('f').
		Default("table").
		EnumVar(&args.outputFormat, "table", "csv", "json")

	extractCli.Arg("extractor-name", "The name of the extractor to use.").
		Required().
		StringVar(&args.extractorName)

	extractCli.Arg("files", "The files to read.").
		Required().
		StringsVar(&args.filePatterns)

	addFileHandlingFlagsTo(globalFlags, extractCli)
}

func (cmd *ExtractCmd) initWithArgs(args *extractArgs, flags *config.GlobalFlags) error {
	resultsWriter, err := resultsWriters.NewExtractionResultsWriter(flags.MultiFileOut,
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

	fileExtractor := ch360.NewFileExtractor(client.Documents, client.Documents, client.Documents)

	cmd.ExtractionService = services.NewParallelExtractionService(fileExtractor, client.Documents,
		progressHandler)
	cmd.ExtractorName = args.extractorName

	return nil
}

// ExecuteExtract is the main entry point for the 'extract' command.
func (cmd *ExtractCmd) Execute(ctx context.Context) error {
	err := cmd.ExtractionService.ExtractAll(ctx, cmd.FilePaths, cmd.ExtractorName)

	return errors.Wrap(err, "extraction failed")
}
