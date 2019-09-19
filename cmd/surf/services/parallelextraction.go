package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/ch360/results"
	"github.com/waives/surf/pool"
	"io"
	"os"
)

//go:generate mockery -name "FileExtractor"

type FileExtractor interface {
	Extract(ctx context.Context, fileContent io.Reader, extractorName string) (*results.ExtractionResult, error)
}

// ParallelExtractionService wraps the ch360.FileExtractor to process multiple files in parallel.
type ParallelExtractionService struct {
	singleFileExtractor    FileExtractor
	documentGetter         ch360.DocumentGetter
	parallelFilesProcessor ParallelFilesProcessor
}

// NewParallelExtractionService constructs a new ParallelExtractionService.
func NewParallelExtractionService(fileExtractor FileExtractor,
	documentGetter ch360.DocumentGetter,
	progressHandler ProgressHandler) *ParallelExtractionService {

	return &ParallelExtractionService{
		singleFileExtractor: fileExtractor,
		documentGetter:      documentGetter,
		parallelFilesProcessor: ParallelFilesProcessor{
			ProgressHandler: progressHandler,
		},
	}
}

func (p *ParallelExtractionService) ExtractAll(ctx context.Context, files []string,
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
				return nil, errors.Wrapf(err, "Error extracting file %s", filename)
			}
			defer file.Close()

			readCloser, err := p.singleFileExtractor.Extract(ctx, file, extractorName)

			return readCloser, errors.Wrapf(err, "Error extracting file %s", filename)
		}
	}

	return p.parallelFilesProcessor.Run(ctx, files, parallelWorkers,
		processorFunc)
}
