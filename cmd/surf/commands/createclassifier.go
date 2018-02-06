package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/docopt/docopt-go"
	"io"
)

//go:generate mockery -name "ClassifierCreator|ClassifierTrainer|ClassifierClient"

const CreateClassifierCommand = "create classifier"

type ClassifierCreator interface {
	Create(name string) error
}

type ClassifierTrainer interface {
	Train(name string, samplesPath string) error
}

type CreateClassifier struct {
	writer         io.Writer
	creator        ClassifierCreator
	deleter        ClassifierDeleter
	trainer        ClassifierTrainer
	classifierName string
	samplesPath    string
}

func NewCreateClassifier(writer io.Writer,
	creator ClassifierCreator,
	trainer ClassifierTrainer,
	deleter ClassifierDeleter,
	classifierName string,
	samplesPath string) *CreateClassifier {
	return &CreateClassifier{
		writer:         writer,
		creator:        creator,
		deleter:        deleter,
		trainer:        trainer,
		classifierName: classifierName,
		samplesPath:    samplesPath,
	}
}

func NewCreateClassifierFromArgs(args docopt.Opts, client *ch360.ApiClient, out io.Writer) (*CreateClassifier, error) {
	var (
		classifierName, _ = args.String("<name>")
		samplesPath, _    = args.String("<samples-zip>")
	)

	return NewCreateClassifier(
		out,
		client.Classifiers,
		client.Classifiers,
		client.Classifiers,
		classifierName,
		samplesPath), nil
}

func (cmd *CreateClassifier) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Creating classifier '%s'... ", cmd.classifierName)

	err := cmd.creator.Create(cmd.classifierName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	fmt.Fprintf(cmd.writer, "Adding samples from file '%s'... ", cmd.samplesPath)
	err = cmd.trainer.Train(cmd.classifierName, cmd.samplesPath)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		cmd.deleter.Delete(cmd.classifierName)
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}
