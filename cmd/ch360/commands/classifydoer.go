package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
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

func (cmd *ClassifyCommand) Execute(filePattern string, classifierName string) error {
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
		return errors.New(fmt.Sprintf("File pattern %s does not match any files", filePattern))
	}

	fmt.Fprintf(cmd.writer, "%-40.40s  %s", "FILE", "DOCUMENT TYPE")
	fmt.Fprintln(cmd.writer)

	for _, filename := range matches {
		documentType, err := cmd.processFile(filename, classifierName)
		if err != nil {
			return errors.New(fmt.Sprintf("Error classifying file %s: %s", filename, err.Error()))
		} else {
			base := filepath.Base(filename)
			fmt.Fprintf(cmd.writer, "%-40.40s  %s", base, documentType)
		}
		fmt.Fprintln(cmd.writer)
	}
	return nil
}

func (cmd *ClassifyCommand) processFile(filePath string, classifierName string) (string, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	documentId, err := cmd.client.CreateDocument(fileContents)
	if err != nil {
		return "", err
	}

	documentType, classifyErr := cmd.client.ClassifyDocument(documentId, classifierName)

	// Always delete the document, even if ClassifyDocument returned an error
	deleteErr := cmd.client.DeleteDocument(documentId)

	if classifyErr != nil {
		return "", classifyErr
	}

	if deleteErr != nil {
		return "", deleteErr
	}

	return documentType, nil
}
