package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/mattn/go-zglob"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

	var jobs []pool.Job
	for _, filename := range matches {
		// The memory of the 'filename' var is reused here, see:
		// https://golang.org/doc/faq#closures_and_goroutines
		// The workaround is to copy it:
		filename := filename // <- copy

		job := pool.NewJob(
			// Performing the work
			func() pool.JobResult {
				result, err := cmd.processFile(ctx, filename, classifierName)

				return pool.JobResult{
					Err:   err,
					Value: result,
				}
			},
			// Handling the result
			func(jr pool.JobResult) {
				if jr.Err != nil {
					fmt.Fprintf(cmd.writer, "Error classifying file %s: %v\n", filename, jr.Err)
				}

				classificationResult := jr.Value.(*types.ClassificationResult)

				if classificationResult != nil {
					fmt.Fprintf(cmd.writer,
						ClassifyOutputFormat,
						filepath.Base(filename),
						classificationResult.DocumentType,
						classificationResult.IsConfident)
				}
			})

		jobs = append(jobs, job)
	}

	workPool := pool.NewPool(jobs, 10)

	// Print results
	fmt.Fprintf(cmd.writer, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	workPool.Run(ctx)

	return nil
}

func (cmd *ClassifyCommand) processFile(ctx context.Context, filePath string, classifierName string) (result *types.ClassificationResult, err error) {
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

	result, err = cmd.client.ClassifyDocument(ctx, documentId, classifierName)

	if documentId != "" {
		// Always delete the document, even if ClassifyDocument returned an error.
		// Don't cancel on ctrl-c.
		err = cmd.client.DeleteDocument(context.Background(), documentId)
	}

	return // named return params
}
