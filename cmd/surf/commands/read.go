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

// ConfigureReadCommand configures kingpin to call ExecuteRead after having successfully parsed
// the cli options.
func ConfigureReadCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	readArgs := &ReadArgs{}

	cmd := app.
		Command("read", "Perform OCR on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			// execute the command
			return ExecuteRead(ctx, readArgs, globalFlags)
		})

	cmd.Flag("format", "The output format. Allowed values: pdf, wvdoc, txt [default: txt].").
		Short('f').
		Default("txt").
		EnumVar(&readArgs.outputFormat, "pdf", "wvdoc", "txt")

	cmd.Arg("files", "The files to read.").
		Required().
		StringsVar(&readArgs.filePatterns)
}

// ExecuteRead is the main entry point for the 'read' command. It has to do a lot of
// setup / instantiation before actually performing OCR on the specified files.
func ExecuteRead(ctx context.Context, readArgs *ReadArgs, globalFlags *config.GlobalFlags) error {

	resultsWriter, err := resultsWriters.NewReaderResultsWriter(globalFlags.MultiFileOut,
		globalFlags.OutputFile, readArgs.outputFormat)

	if err != nil {
		return err
	}

	progressHandler := progress.NewProgressHandler(resultsWriter,
		globalFlags.ShowProgress, os.Stderr)

	readMode := readModes[readArgs.outputFormat]

	// ensure we're not printing binary data to the console
	if !config.IsOutputRedirected() &&
		readMode.IsBinary() &&
		!globalFlags.IsOutputSpecified() {
		return errors.New("you must use '-o' or '-m' or redirect stdout when the output " +
			"file format is pdf or wvdoc")
	}

	filePaths, err := GlobMany(readArgs.filePatterns)
	if err != nil {
		return err
	}

	client, err := initApiClient(globalFlags.ClientId,
		globalFlags.ClientSecret,
		globalFlags.LogHttp)

	if err != nil {
		return err
	}

	fileReader := ch360.NewFileReader(client.Documents, client.Documents, client.Documents)

	return services.NewParallelReaderService(fileReader, client.Documents, progressHandler).
		ReadAll(ctx, filePaths, readMode)
}

var readModes = map[string]ch360.ReadMode{
	"wvdoc": ch360.ReadWvdoc,
	"pdf":   ch360.ReadPDF,
	"txt":   ch360.ReadText,
}
