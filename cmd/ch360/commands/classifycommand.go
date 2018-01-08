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

var ClassifyOutputFormat = "%-44.44s %-24.24s %v"

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
	fmt.Fprintln(cmd.writer)

	for _, filename := range matches {
		result, err := cmd.processFile(ctx, filename, classifierName)
		if err != nil {
			return errors.New(fmt.Sprintf("Error classifying file %s: %s", filename, err.Error()))
		} else {
			base := filepath.Base(filename)
			fmt.Fprintf(cmd.writer, ClassifyOutputFormat, base, result.DocumentType, result.IsConfident)
		}
		fmt.Fprintln(cmd.writer)
	}
	return nil
}

func (cmd *ClassifyCommand) processFile(ctx context.Context, filePath string, classifierName string) (*types.ClassificationResult, error) {

	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var (
		result     *types.ClassificationResult
		documentId string
	)

	errChan := make(chan error, 1)
	go func() {
		errChan <- func() error {
			documentId, err = cmd.client.CreateDocument(ctx, fileContents)
			if err != nil {
				return err
			}

			result, err = cmd.client.ClassifyDocument(ctx, documentId, classifierName)

			if err != nil {
				return err
			}

			return nil
		}()
	}()

	var classifyErr error
	var deleteErr error
	cancelled := false
	select {
	case <-ctx.Done():
		cancelled = true
	case classifyErr = <-errChan:
	}

	if documentId != "" {
		// Always delete the document, even if ClassifyDocument returned an error.
		// Don't cancel on ctrl-c
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
