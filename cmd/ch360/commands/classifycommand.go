package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/mattn/go-zglob"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type ClassifyCommand struct {
	writer io.Writer
	client ch360.DocumentCreatorDeleterClassifier
}

func NewClassifyCommand(writer io.Writer, client ch360.DocumentCreatorDeleterClassifier) *ClassifyCommand {
	return &ClassifyCommand{
		writer: writer,
		client: client,
	}
}

var ClassifyOutputFormat = "%-36.36s %-32.32s %v\n"

type job struct {
	filename       string
	classifierName string
}

type jobResult struct {
	err    error
	result *types.ClassificationResult
	job    *job
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

	fileCount := len(matches)
	if fileCount == 0 {
		return errors.New(fmt.Sprintf("File glob pattern %s does not match any files. Run 'ch360 -h' for glob pattern examples.", filePattern))
	}

	// Set up jobs channel
	jobsChan := make(chan job, 0)
	go func() {
		for _, filename := range matches {
			select {
			case <-ctx.Done():
				break
			default:
				jobsChan <- job{
					classifierName: classifierName,
					filename:       filename,
				}
			}
		}
		close(jobsChan)
	}()

	// Set up results channels
	resultsChan := make(chan jobResult, 0)

	// Start processing in background
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			for job := range jobsChan {
				result, err := cmd.processFile(ctx, job.filename, job.classifierName)
				resultsChan <- jobResult{
					err:    err,
					result: result,
					job:    &job,
				}
			}
			wg.Done()
		}(i)
	}

	go func() {
		// Wait for all workers to complete, then close the results channel
		wg.Wait()
		close(resultsChan)
	}()

	// Print results
	fmt.Fprintf(cmd.writer, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")

	for jr := range resultsChan {
		err, result := jr.err, jr.result

		if err == context.Canceled {
			continue
		}

		if err != nil {
			fmt.Fprintf(cmd.writer, "Error classifying file %s: %s\n", jr.job.filename, err.Error())
		} else {
			fmt.Fprintf(cmd.writer, ClassifyOutputFormat, filepath.Base(jr.job.filename), result.DocumentType, result.IsConfident)
		}

	}
	return nil
}

func (cmd *ClassifyCommand) processFile(ctx context.Context, filePath string, classifierName string) (*types.ClassificationResult, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	documentId, err := cmd.client.CreateDocument(context.Background(), fileContents)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)
	var result *types.ClassificationResult
	go func() {
		result, err = cmd.client.ClassifyDocument(ctx, documentId, classifierName)

		errChan <- err
	}()

	var classifyErr error
	var deleteErr error

	var cancelled = false
	select {
	case <-ctx.Done():
		cancelled = true
	case classifyErr = <-errChan:
	}

	if documentId != "" {
		// Always delete the document, even if ClassifyDocument returned an error.
		// Don't cancel on ctrl-c.
		deleteErr = cmd.client.DeleteDocument(context.Background(), documentId)
	}

	if classifyErr != nil {
		return nil, classifyErr
	}

	if deleteErr != nil {
		return nil, deleteErr
	}

	if cancelled {
		return nil, ctx.Err()
	}

	return result, nil
}
