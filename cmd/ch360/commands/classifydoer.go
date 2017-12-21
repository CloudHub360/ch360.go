package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
	"io/ioutil"
	"os"
)

type ClassifyDoer struct {
	writer io.Writer
	client ch360.DocumentCreatorDeleterClassifier //TODO: Should this be in ch360 or not?
}

func NewClassifyDoer(writer io.Writer, client ch360.DocumentCreatorDeleterClassifier) *ClassifyDoer {
	return &ClassifyDoer{
		writer: writer,
		client: client,
	}
}

func (cmd *ClassifyDoer) Execute(filePath, classifierName string) error {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			//TODO: Revisit text of this when accepting patterns for multiple files
			return errors.New(fmt.Sprintf("File %s does not exist", filePath))
		} else {
			return err
		}
	}

	documentId, err := cmd.client.CreateDocument(fileContents)
	if err != nil {
		return err
	}

	documentType, classifyErr := cmd.client.ClassifyDocument(documentId, classifierName)

	// Always delete the document, even if ClassifyDocument returned an error
	deleteErr := cmd.client.DeleteDocument(documentId)

	if classifyErr != nil {
		return classifyErr
	}

	if deleteErr != nil {
		return deleteErr
	}

	fmt.Fprintln(cmd.writer, documentType)
	return nil
}
