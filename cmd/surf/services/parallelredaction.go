package services

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/pkg/errors"
	"io"
	"os"
)

//go:generate mockery -name "FileRedactor"

type FileRedactor interface {
	Redact(ctx context.Context, fileContent io.Reader, extractorName string) (io.ReadCloser, error)
}

// ParallelRedactionService wraps the ch360.FileRedactor to process multiple files in parallel.
type ParallelRedactionService struct {
	singleFileRedactor     FileRedactor
	documentGetter         ch360.DocumentGetter
	parallelFilesProcessor ParallelFilesProcessor
}

// NewParallelRedactionService constructs a new ParallelRedactionService.
func NewParallelRedactionService(fileRedactor FileRedactor,
	documentGetter ch360.DocumentGetter,
	progressHandler ProgressHandler) *ParallelRedactionService {

	return &ParallelRedactionService{
		singleFileRedactor: fileRedactor,
		documentGetter:     documentGetter,
		parallelFilesProcessor: ParallelFilesProcessor{
			ProgressHandler: progressHandler,
		},
	}
}

func (p *ParallelRedactionService) RedactAllWithExtractor(ctx context.Context, files []string,
	extractorName string) error {

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
				return nil, errors.Wrapf(err, "Error redacting file %s", filename)
			}
			defer file.Close()

			readCloser, err := p.singleFileRedactor.Redact(ctx, file, extractorName)

			return readCloser, errors.Wrapf(err, "Error redacting file %s", filename)
		}
	}

	return p.parallelFilesProcessor.Run(ctx, files, parallelWorkers,
		processorFunc)
}
