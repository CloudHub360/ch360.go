package services

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/pkg/errors"
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

	// Get the current number of documents, so we know how many slots are available
	docs, err := p.documentGetter.GetAll(ctx)
	if err != nil {
		return err
	}

	// Limit the number of workers to the number of available doc slots
	parallelWorkers := ch360.TotalDocumentSlots - len(docs)

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
