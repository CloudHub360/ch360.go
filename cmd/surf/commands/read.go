package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"io"
	"os"

	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/progress"
	"github.com/CloudHub360/ch360.go/pool"
	"gopkg.in/alecthomas/kingpin.v2"
)

const ReadFilesCommand = "read"

//go:generate mockery -name "FileReader|FilesProcessor"

type FileReader interface {
	Read(ctx context.Context, fileContent io.Reader, mode ch360.ReadMode) (io.ReadCloser, error)
}

type FilesProcessor interface {
	RunWithGlob(ctx context.Context,
		filePatterns []string,
		parallelism int,
		processorFuncFactory ProcessorFuncFactory) error
}

type ReadCmd struct {
	fileReader FileReader
	//parallelWorkers int
	//filesProcessor  FilesProcessor
	outputFormat string
	filePatterns []string
	globalFlags  *config.GlobalFlags
}

func ConfigureReadCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	readCmd := &ReadCmd{
		globalFlags: globalFlags,
	}
	cmd := app.
		Command("read", "Perform OCR on a file or set of files.").
		Action(func(parseContext *kingpin.ParseContext) error {
			return readCmd.Execute(ctx)
		})

	cmd.Arg("format", "The output format. Allowed values: pdf, wvdoc, txt.").
		Required().
		EnumVar(&readCmd.outputFormat, "pdf", "wvdoc", "txt")

	cmd.Arg("files", "The files to read.").
		Required().
		StringsVar(&readCmd.filePatterns)
}

func (cmd *ReadCmd) processorFor(ctx context.Context, filename string,
	readMode ch360.ReadMode) pool.ProcessorFunc {
	return func() (interface{}, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, wrapErr(filename, err)
		}

		readCloser, err := cmd.fileReader.Read(ctx, file, readMode)

		return readCloser, wrapErr(filename, err)
	}
}

func wrapErr(filename string, err error) error {
	if err == nil {
		return nil
	}
	return errors.New(fmt.Sprintf("Error reading file %s: %v", filename, err))
}

func (cmd *ReadCmd) Execute(ctx context.Context) error {
	client, err := initApiClient(cmd.globalFlags.ClientId, cmd.globalFlags.ClientSecret, cmd.globalFlags.LogHttp)

	if err != nil {
		return err
	}

	fileExtension := fmt.Sprintf(".ocr.%s", cmd.outputFormat)

	resultsWriter, err := resultsWriters.NewResultsWriterFor(cmd.globalFlags, fileExtension,
		config.Read)

	progressHandler := progress.NewProgressHandler(resultsWriter,
		cmd.globalFlags.ShowProgress, os.Stderr)

	if err != nil {
		return err
	}

	filesProcessor := &ParallelFilesProcessor{ProgressHandler: progressHandler}

	var readMode = ch360.ReadText
	switch cmd.outputFormat {
	case "wvdoc":
		readMode = ch360.ReadWvdoc
	case "pdf":
		readMode = ch360.ReadPDF
	}

	// ensure we're not printing binary data to the console
	if !config.IsOutputRedirected() &&
		readMode.IsBinary() &&
		!cmd.globalFlags.IsOutputSpecified() {
		return errors.New("you must use '-o' or '-m' or redirect stdout when the output " +
			"file format is pdf")
	}

	cmd.fileReader = ch360.NewFileReader(client.Documents, client.Documents, client.Documents)

	// Get the current number of documents, so we know how many slots are available
	docs, err := client.Documents.GetAll(ctx)
	if err != nil {
		return err
	}

	// Limit the number of workers to the number of available doc slots
	parallelWorkers := min(10, ch360.TotalDocumentSlots-len(docs))

	return filesProcessor.RunWithGlob(ctx,
		cmd.filePatterns,
		parallelWorkers,
		func(ctx context.Context, filename string) pool.ProcessorFunc {
			return cmd.processorFor(ctx, filename, readMode)
		})
}

//
//func NewReadFilesCommand(reader FileReader,
//	processor FilesProcessor,
//	mode ch360.ReadMode,
//	filePatterns []string,
//	parallelWorkers int,
//	getter ch360.DocumentGetter) *ReadCmd {
//	return &ReadCmd{
//		filesProcessor:  processor,
//		fileReader:      reader,
//		documentGetter:  getter,
//		parallelWorkers: parallelWorkers,
//		filePatterns:    filePatterns,
//		mode:            mode,
//	}
//}
//
//func NewReadFilesCommandFromArgs(params *config.GlobalFlags, client *ch360.ApiClient, readFormat string, filePatterns []string) (*ReadCmd, error) {
//
//	progressHandler, err := progress.NewProgressHandlerFor(params, os.Stderr)
//
//	if err != nil {
//		return nil, err
//	}
//
//	filesProcessor := &ParallelFilesProcessor{
//		ProgressHandler: progressHandler,
//	}
//
//	var readMode = ch360.ReadText
//	switch readFormat {
//	case "wvdoc":
//		readMode = ch360.ReadWvdoc
//	case "pdf":
//		readMode = ch360.ReadPDF
//	}
//
//	// ensure we're not printing binary data to the console
//	if !config.IsOutputRedirected() &&
//		readMode.IsBinary() &&
//		!params.IsOutputSpecified() {
//		return nil, errors.New("You must use '-o' or '-m' or redirect stdout when the output " +
//			"file format is pdf.")
//	}
//
//	fileReader := ch360.NewFileReader(client.Documents, client.Documents, client.Documents)
//	return NewReadFilesCommand(fileReader,
//		filesProcessor,
//		readMode,
//		filePatterns,
//		10,
//		client.Documents), nil
//
//}

func (cmd ReadCmd) Usage() string {
	return ReadFilesCommand
}
