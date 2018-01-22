package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/mattn/go-zglob"
	"io/ioutil"
	"os"
)

type ClassifyCommand struct {
	documentClassifier ch360.DocumentClassifier
	documentCreator    ch360.DocumentCreator
	documentDeleter    ch360.DocumentDeleter
	documentGetter     ch360.DocumentGetter
	parallelWorkers    int
	progressHandler    ClassifyProgressHandler
}

//go:generate mockery -name ClassifyProgressHandler
type ClassifyProgressHandler interface {
	Notify(filename string, result *types.ClassificationResult) error
	NotifyErr(filename string, err error)
	NotifyStart(totalJobs int) error
	NotifyFinish() error
}

func NewClassifyCommand(progressHandler ClassifyProgressHandler,
	docClassifier ch360.DocumentClassifier,
	docCreator ch360.DocumentCreator,
	docDeleter ch360.DocumentDeleter,
	docGetter ch360.DocumentGetter,
	parallelism int) *ClassifyCommand {
	return &ClassifyCommand{
		progressHandler:    progressHandler,
		documentClassifier: docClassifier,
		documentCreator:    docCreator,
		documentDeleter:    docDeleter,
		documentGetter:     docGetter,
		parallelWorkers:    parallelism,
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (cmd *ClassifyCommand) handlerFor(cancel context.CancelFunc, filename string, errs *[]error) pool.HandlerFunc {
	return func(value interface{}, err error) {
		if err != nil {
			errMsg := fmt.Sprintf("Error classifying file %s: %v", filename, err)
			err = errors.New(errMsg)
			cmd.progressHandler.NotifyErr(filename, err)
			*errs = append(*errs, err)
			// Don't process any more if there's an error
			cancel()
		} else {
			classificationResult := value.(*types.ClassificationResult)

			if err = cmd.progressHandler.Notify(filename, classificationResult); err != nil {
				cancel()
			}
		}
	}
}

func (cmd *ClassifyCommand) Execute(ctx context.Context, filePattern string, classifierName string) error {
	matches, err := zglob.Glob(filePattern)
	if err != nil {
		if os.IsNotExist(err) {
			// The file pattern is for a specific (single) file that doesn't exist
			return errors.New(fmt.Sprintf("File %s does not exist", filePattern))
		} else {
			return err
		}
	}

	// Get the current number of documents, so we know how many slots are available
	docs, err := cmd.documentGetter.GetAll(ctx)
	if err != nil {
		return err
	}
	// Limit the number of workers to the number of available doc slots
	cmd.parallelWorkers = min(cmd.parallelWorkers, ch360.TotalDocumentSlots-len(docs))

	fileCount := len(matches)
	if fileCount == 0 {
		return errors.New(fmt.Sprintf("File glob pattern %s does not match any files. Run 'surf -h' for glob pattern examples.", filePattern))
	}

	ctx, cancel := context.WithCancel(ctx)

	var (
		processFileJobs []pool.Job
		errors          []error
	)

	for _, filename := range matches {
		// The memory of the 'filename' var is reused here, see:
		// https://golang.org/doc/faq#closures_and_goroutines
		// The workaround is to copy it:
		filename := filename // <- copy

		processFileJob := pool.NewJob(
			func() (interface{}, error) {
				return cmd.processFile(ctx, filename, classifierName)
			},
			cmd.handlerFor(cancel, filename, &errors))

		processFileJobs = append(processFileJobs, processFileJob)
	}

	workPool := pool.NewPool(processFileJobs, cmd.parallelWorkers)

	// Print results
	cmd.progressHandler.NotifyStart(len(processFileJobs))
	defer cmd.progressHandler.NotifyFinish()
	workPool.Run(ctx)

	// Just return the first error.
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

func (cmd *ClassifyCommand) processFile(ctx context.Context, filePath string, classifierName string) (result *types.ClassificationResult, err error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	documentId, err := cmd.documentCreator.Create(context.Background(), fileContents)
	if err != nil {
		return nil, err
	}

	result, classifyErr := cmd.documentClassifier.Classify(ctx, documentId, classifierName)

	if documentId != "" {
		// Always delete the document, even if Classify returned an error.
		// Don't cancel on ctrl-c.
		err = cmd.documentDeleter.Delete(context.Background(), documentId)
	}

	// Return the classify err if we have one
	if classifyErr != nil {
		err = classifyErr
	}

	return // named return params
}
