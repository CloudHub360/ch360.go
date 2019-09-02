package commands

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/services"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/progress"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type redactWithExtractorArgs struct {
	extractorName string
	filePatterns  []string
}

type RedactWithExtractorCmd struct {
	ExtractorName    string
	FilePaths        []string
	RedactionService RedactWithExtractorService
}

//go:generate mockery -name "RedactWithExtractorService"
type RedactWithExtractorService interface {
	RedactAllWithExtractor(ctx context.Context, files []string, extractorName string) error
}

func ConfigureRedactWithExtractionCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	args := &redactWithExtractorArgs{}
	cmd := &RedactWithExtractorCmd{}
	redactCli := app.
		Command("redact", "Perform data redaction on a file or set of files.")

	redactWithExtractorCli := redactCli.Command("with-extractor",
		"Use fields from an extractor to define areas to redact. ").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := cmd.initWithArgs(args, globalFlags)
			if err != nil {
				return err
			}
			return cmd.Execute(ctx)
		})

	redactWithExtractorCli.Arg("extractor-name", "The name of the extractor to use.").
		Required().
		StringVar(&args.extractorName)

	redactWithExtractorCli.Arg("files", "The files to read.").
		Required().
		StringsVar(&args.filePatterns)
}

func (cmd *RedactWithExtractorCmd) initWithArgs(args *redactWithExtractorArgs, flags *config.GlobalFlags) error {
	resultsWriter, err := resultsWriters.NewRedactResultsWriter(flags.MultiFileOut, flags.OutputFile)

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

	fileRedactor := ch360.NewFileRedactor(client.Documents, client.Documents, client.Documents,
		client.Documents)

	cmd.RedactionService = services.NewParallelRedactionService(fileRedactor, client.Documents, progressHandler)
	cmd.ExtractorName = args.extractorName

	return nil
}

// ExecuteRedact is the main entry point for the 'redact' command.
func (cmd *RedactWithExtractorCmd) Execute(ctx context.Context) error {
	err := cmd.RedactionService.RedactAllWithExtractor(ctx, cmd.FilePaths, cmd.ExtractorName)

	return errors.Wrap(err, "Redaction failed")
}
