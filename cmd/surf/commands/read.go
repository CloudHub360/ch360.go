package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/progress"
	"github.com/CloudHub360/ch360.go/pool"
	"io"
	"os"
)

const ReadFilesCommand = "read"

//go:generate mockery -name "FileReader"

type FileReader interface {
	Read(ctx context.Context, fileContent io.Reader, mode ch360.ReadMode) (io.ReadCloser, error)
}

type FilesProcessor interface {
	RunWithGlob(ctx context.Context,
		filesPattern string,
		parallelism int,
		processorFuncFactory ProcessorFuncFactory) error
}

type Read struct {
	fileReader      FileReader
	documentGetter  ch360.DocumentGetter
	parallelWorkers int
	filesProcessor  FilesProcessor
	mode            ch360.ReadMode
	filesPattern    string
}

func (cmd *Read) ProcessorFor(ctx context.Context, filename string) pool.ProcessorFunc {
	return func() (interface{}, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, wrapErr(filename, err)
		}

		readCloser, err := cmd.fileReader.Read(ctx, file, cmd.mode)

		return readCloser, wrapErr(filename, err)
	}
}

func wrapErr(filename string, err error) error {
	if err == nil {
		return nil
	}
	return errors.New(fmt.Sprintf("Error reading file %s: %v", filename, err))
}

func (cmd *Read) Execute(ctx context.Context) error {

	// Get the current number of documents, so we know how many slots are available
	docs, err := cmd.documentGetter.GetAll(ctx)
	if err != nil {
		return err
	}

	// Limit the number of workers to the number of available doc slots
	parallelWorkers := min(cmd.parallelWorkers, ch360.TotalDocumentSlots-len(docs))

	return cmd.filesProcessor.RunWithGlob(ctx,
		cmd.filesPattern,
		parallelWorkers,
		cmd,
		cmd)
}

func NewReadFilesCommandFromArgs(params *config.RunParams, client *ch360.ApiClient) (*Read, error) {

	progressHandler, err := progress.NewProgressHandlerFor(params, os.Stderr)

	if err != nil {
		return nil, err
	}

	filesProcessor := &ParallelFilesProcessor{
		ProgressHandler: progressHandler,
	}

	var readMode = ch360.ReadText
	if params.ReadPDF {
		readMode = ch360.ReadPDF
	} else if params.ReadWvdoc {
		readMode = ch360.ReadWvdoc
	}

	return &Read{
		filesProcessor:  filesProcessor,
		fileReader:      ch360.NewFileReader(client.Documents, client.Documents, client.Documents),
		documentGetter:  client.Documents,
		parallelWorkers: 10,
		filesPattern:    params.FilePattern,
		mode:            readMode,
	}, nil
}

func (cmd Read) Usage() string {
	return ReadFilesCommand
}
