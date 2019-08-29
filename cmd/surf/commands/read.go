package commands

import (
	"context"
	"github.com/CloudHub360/ch360.go/cmd/surf/services"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"github.com/pkg/errors"
	"os"

	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/progress"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ReadArgs struct {
	outputFormat string
	filePatterns []string
}

//go:generate mockery -name ReaderService
type ReaderService interface {
	ReadAll(ctx context.Context, files []string, readMode ch360.ReadMode) error
}

// ReadCmd represents the 'read' command. It relies on a 'ReaderService' to perform the actual OCR.
type ReadCmd struct {
	FilePaths     []string
	ReaderService ReaderService
	ReadMode      ch360.ReadMode
}

func (cmd *ReadCmd) initFromArgs(args *ReadArgs, globalFlags *config.GlobalFlags) error {
	resultsWriter, err := resultsWriters.NewReaderResultsWriter(globalFlags.MultiFileOut,
		globalFlags.OutputFile, args.outputFormat)

	if err != nil {
		return err
	}

	progressHandler := progress.NewProgressHandler(resultsWriter, globalFlags.ShowProgress, os.Stderr)

	cmd.ReadMode = readModes[args.outputFormat]

	// ensure we're not printing binary data to the console
	if !config.IsOutputRedirected() &&
		cmd.ReadMode.IsBinary() &&
		!globalFlags.IsOutputSpecified() {
		return errors.New("you must use '-o' or '-m' or redirect stdout when the output " +
			"file format is pdf or wvdoc")
	}

	cmd.FilePaths, err = GlobMany(args.filePatterns)
	if err != nil {
		return err
	}

	client, err := initApiClient(globalFlags.ClientId,
		globalFlags.ClientSecret,
		globalFlags.LogHttp)

	if err != nil {
		return err
	}

	singleFileReader := ch360.NewFileReader(client.Documents, client.Documents, client.Documents)

	cmd.ReaderService = services.NewParallelReaderService(singleFileReader, client.Documents,
		progressHandler)

	return nil
}

// ConfigureReadCommand configures kingpin to call ExecuteRead after having successfully parsed
// the cli options.
func ConfigureReadCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	readArgs := &ReadArgs{}

	readCmd := &ReadCmd{}

	cliCmd := app.
		Command("read", "Perform OCR on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := readCmd.initFromArgs(readArgs, globalFlags)
			if err != nil {
				return err
			}
			return readCmd.Execute(ctx)
		})

	cliCmd.Flag("format", "The output format. Allowed values: pdf, wvdoc, txt [default: txt].").
		Short('f').
		Default("txt").
		EnumVar(&readArgs.outputFormat, "pdf", "wvdoc", "txt")

	cliCmd.Arg("files", "The files to read.").
		Required().
		StringsVar(&readArgs.filePatterns)
}

// Execute is the main entry point for the 'read' command.
func (cmd *ReadCmd) Execute(ctx context.Context) error {
	err := cmd.ReaderService.ReadAll(ctx, cmd.FilePaths, cmd.ReadMode)
	return errors.Wrap(err, "Read failed")
}

var readModes = map[string]ch360.ReadMode{
	"wvdoc": ch360.ReadWvdoc,
	"pdf":   ch360.ReadPDF,
	"txt":   ch360.ReadText,
}
