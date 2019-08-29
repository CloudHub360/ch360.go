package services

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/pkg/errors"
	"io"
	"os"
)

//go:generate mockery -name "FileReader|FileGlobProcessor"

type FileReader interface {
	Read(ctx context.Context, fileContent io.Reader, mode ch360.ReadMode) (io.ReadCloser, error)
}

type FileGlobProcessor interface {
	RunWithGlob(ctx context.Context,
		filePatterns []string,
		parallelism int,
		processorFuncFactory ProcessorFuncFactory) error
}

// ParallelReaderService wraps the ch360.FileReader to process multiple files in parallel.
type ParallelReaderService struct {
	singleFileReader       FileReader
	documentGetter         ch360.DocumentGetter
	parallelFilesProcessor ParallelFilesProcessor
}

// NewParallelReaderService constructs a new ParallelReaderService.
func NewParallelReaderService(fileReader FileReader,
	documentGetter ch360.DocumentGetter,
	progressHandler ProgressHandler) *ParallelReaderService {

	return &ParallelReaderService{
		singleFileReader: fileReader,
		documentGetter:   documentGetter,
		parallelFilesProcessor: ParallelFilesProcessor{
			ProgressHandler: progressHandler,
		},
	}
}

func (p *ParallelReaderService) ReadAll(ctx context.Context, files []string,
	readMode ch360.ReadMode) error {

	// Limit the number of workers to the number of available doc slots
	parallelWorkers, err := ch360.GetFreeDocSlots(ctx, p.documentGetter, ch360.TotalDocumentSlots)

	if err != nil {
		return err
	}

	// called in parallel, once per file
	processorFunc := func(ctx context.Context, filename string) pool.ProcessorFunc {
		return func() (interface{}, error) {
			file, err := os.Open(filename)
			if err != nil {
				return nil, errors.Wrapf(err, "Error reading file %s", filename)
			}
			defer file.Close()

			readCloser, err := p.singleFileReader.Read(ctx, file, readMode)

			return readCloser, errors.Wrapf(err, "Error reading file %s", filename)
		}
	}

	return p.parallelFilesProcessor.Run(ctx, files, parallelWorkers,
		processorFunc)
}
