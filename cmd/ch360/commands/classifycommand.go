package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/mattn/go-zglob"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
)

type ClassifyCommand struct {
	writer    io.Writer
	client    ch360.DocumentCreatorDeleterClassifier
	interrupt chan os.Signal
}

func NewClassifyCommand(writer io.Writer, client ch360.DocumentCreatorDeleterClassifier) *ClassifyCommand {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt)

	// Block until a signal is received.
	//<-quitChan

	return &ClassifyCommand{
		writer:    writer,
		client:    client,
		interrupt: quitChan,
	}
}

var ClassifyOutputFormat = "%-44.44s %-24.24s %v"

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
		return errors.New(fmt.Sprintf("File glob pattern %s does not match any files. Run 'ch360 -h' for glob pattern examples.", filePattern))
	}

	fmt.Fprintf(cmd.writer, ClassifyOutputFormat, "FILE", "DOCUMENT TYPE", "CONFIDENT")
	fmt.Fprintln(cmd.writer)

	for _, filename := range matches {
		result, err := cmd.processFile(filename, classifierName)
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

func (cmd *ClassifyCommand) processFile(filePath string, classifierName string) (*types.ClassificationResult, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	documentId, err := cmd.client.CreateDocument(fileContents)
	if err != nil {
		return nil, err
	}

	go func() {
		<-cmd.interrupt
		cmd.client.DeleteDocument(documentId)
	}()

	result, classifyErr := cmd.client.ClassifyDocument(documentId, classifierName)

	// Always delete the document, even if ClassifyDocument returned an error
	deleteErr := cmd.client.DeleteDocument(documentId)

	if classifyErr != nil {
		return nil, classifyErr
	}

	if deleteErr != nil {
		return nil, deleteErr
	}

	return result, nil
}
