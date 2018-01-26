package commands

import (
	"fmt"
	"io"
)

//go:generate mockery -name "ClassifierCreator|ClassifierTrainer|ClassifierClient"

type ClassifierCreator interface {
	Create(name string) error
}

type ClassifierTrainer interface {
	Train(name string, samplesPath string) error
}

type CreateClassifier struct {
	writer  io.Writer
	creator ClassifierCreator
	deleter ClassifierDeleter
	trainer ClassifierTrainer
}

func NewCreateClassifier(writer io.Writer,
	creator ClassifierCreator,
	trainer ClassifierTrainer,
	deleter ClassifierDeleter) *CreateClassifier {
	return &CreateClassifier{
		writer:  writer,
		creator: creator,
		deleter: deleter,
		trainer: trainer,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string, samplesPath string) error {
	fmt.Fprintf(cmd.writer, "Creating classifier '%s'... ", classifierName)

	err := cmd.creator.Create(classifierName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	fmt.Fprintf(cmd.writer, "Adding samples from file '%s'... ", samplesPath)
	err = cmd.trainer.Train(classifierName, samplesPath)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		cmd.deleter.Delete(classifierName)
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}
