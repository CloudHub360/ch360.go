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

var ClassifyOutputFormat = "%-44.44s %-24.24s %v\n"

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

	if len(matches) == 0 {
		return errors.New(fmt.Sprintf("File glob pattern %s does not match any files. Run 'ch360 -h' for glob pattern examples.", filePattern))
	}

	fmt.Fprintf(cmd.writer, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")

	for _, filename := range matches {
		result, err := cmd.processFile(ctx, filename, classifierName)
		if err != nil {
			return errors.New(fmt.Sprintf("Error classifying file %s: %s", filename, err.Error()))
		} else if result != nil {
			base := filepath.Base(filename)
			fmt.Fprintf(cmd.writer, ClassifyOutputFormat, base, result.DocumentType, result.IsConfident)
		}

		if ctx.Err() == context.Canceled {
			return nil
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

		if err != nil {
			errChan <- err
		}
	}()

	var classifyErr error
	var deleteErr error

	select {
	case <-ctx.Done():
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

	return result, nil
}
