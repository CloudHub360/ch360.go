package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/pkg/errors"
	"io"
	"os"
)

//go:generate mockery -name "ClassifierCreator|ClassifierTrainer|ClassifierClient"

const CreateClassifierCommand = "create classifier"

type ClassifierCreator interface {
	Create(ctx context.Context, name string) error
}

type ClassifierTrainer interface {
	Train(ctx context.Context, name string, samplesArchive io.Reader) error
}

type CreateClassifier struct {
	writer         io.Writer
	creator        ClassifierCreator
	deleter        ClassifierDeleter
	trainer        ClassifierTrainer
	classifierName string
	samplesArchive io.ReadCloser
}

func NewCreateClassifier(writer io.Writer,
	creator ClassifierCreator,
	trainer ClassifierTrainer,
	deleter ClassifierDeleter,
	classifierName string,
	samplesArchive io.ReadCloser) *CreateClassifier {
	return &CreateClassifier{
		writer:         writer,
		creator:        creator,
		deleter:        deleter,
		trainer:        trainer,
		classifierName: classifierName,
		samplesArchive: samplesArchive,
	}
}

func NewCreateClassifierFromArgs(params *config.RunParams, client *ch360.ApiClient, out io.Writer) (*CreateClassifier, error) {
	samplesArchive, err := os.Open(params.SamplesPath)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be found.", params.SamplesPath))
	}

	return NewCreateClassifier(
		out,
		client.Classifiers,
		client.Classifiers,
		client.Classifiers,
		params.Name,
		samplesArchive), nil
}

func (cmd *CreateClassifier) Execute(ctx context.Context) error {
	defer cmd.samplesArchive.Close()

	fmt.Fprintf(cmd.writer, "Creating classifier '%s'... ", cmd.classifierName)

	err := cmd.creator.Create(ctx, cmd.classifierName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	fmt.Fprintf(cmd.writer, "Adding samples... ")
	err = cmd.trainer.Train(ctx, cmd.classifierName, cmd.samplesArchive)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		cmd.deleter.Delete(ctx, cmd.classifierName)
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}

func (cmd CreateClassifier) Usage() string {
	return CreateClassifierCommand
}
