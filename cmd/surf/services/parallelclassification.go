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

//go:generate mockery -name "FileClassifier"

type FileClassifier interface {
	Classify(ctx context.Context, fileContent io.Reader, classifierName string) (*results.ClassificationResult, error)
}

// ParallelClassificationService wraps the ch360.FileClassifier to process multiple files in parallel.
type ParallelClassificationService struct {
	singleFileClassifier   FileClassifier
	documentGetter         ch360.DocumentGetter
	parallelFilesProcessor ParallelFilesProcessor
}

// NewParallelClassificationService constructs a new ParallelClassificationService.
func NewParallelClassificationService(fileClassifier FileClassifier,
	documentGetter ch360.DocumentGetter,
	progressHandler ProgressHandler) *ParallelClassificationService {

	return &ParallelClassificationService{
		singleFileClassifier: fileClassifier,
		documentGetter:       documentGetter,
		parallelFilesProcessor: ParallelFilesProcessor{
			ProgressHandler: progressHandler,
		},
	}
}

func (p *ParallelClassificationService) ClassifyAll(ctx context.Context, files []string,
	classifierName string) error {

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
				return nil, errors.Wrapf(err, "Error classifying file %s", filename)
			}
			defer file.Close()

			readCloser, err := p.singleFileClassifier.Classify(ctx, file, classifierName)

			return readCloser, errors.Wrapf(err, "Error classifying file %s", filename)
		}
	}

	return p.parallelFilesProcessor.Run(ctx, files, parallelWorkers,
		processorFunc)
}
