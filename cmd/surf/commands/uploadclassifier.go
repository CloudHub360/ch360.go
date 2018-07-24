package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"io"
	"os"
)

//go:generate mockery -name "ClassifierUploader"

const UploadClassifierCommand = "upload classifier"

type ClassifierUploader interface {
	Upload(ctx context.Context, name string, contents io.Reader) error
}

type UploadClassifier struct {
	writer         io.Writer
	uploader       ClassifierUploader
	deleter        ClassifierDeleter
	classifierName string
	classifierFile io.ReadCloser
}

func NewUploadClassifier(writer io.Writer,
	uploader ClassifierUploader,
	classifierName string,
	classifierFile io.ReadCloser) *UploadClassifier {
	return &UploadClassifier{
		writer:         writer,
		uploader:       uploader,
		classifierName: classifierName,
		classifierFile: classifierFile,
	}
}

func NewUploadClassifierFromArgs(params *config.RunParams, client *ch360.ApiClient, out io.Writer) (*UploadClassifier, error) {
	classifierFile, err := os.Open(params.ClassifierPath)

	if err != nil {
		return nil, err
	}

	return NewUploadClassifier(
		out,
		client.Classifiers,
		params.Name,
		classifierFile), nil
}

func (cmd *UploadClassifier) Execute(ctx context.Context) error {
	defer cmd.classifierFile.Close()

	fmt.Fprintf(cmd.writer, "Creating classifier '%s'... ", cmd.classifierName)

	err := cmd.uploader.Upload(ctx, cmd.classifierName, cmd.classifierFile)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}

func (cmd UploadClassifier) Usage() string {
	return UploadClassifierCommand
}
