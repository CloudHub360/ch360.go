package commands

import (
	"fmt"
	"io"
)

//go:generate mockery -name "Creator|Trainer|CreatorTrainer"

type Creator interface {
	Create(name string) error
}

type Trainer interface {
	Train(name string, samplesPath string) error
}

type CreatorTrainer interface {
	Creator
	Trainer
}

type CreateClassifier struct {
	writer           io.Writer
	client           CreatorTrainer
	deleteClassifier ClassifierCommand
}

func NewCreateClassifier(writer io.Writer, client CreatorTrainer, deleteClassifier ClassifierCommand) *CreateClassifier {
	return &CreateClassifier{
		writer:           writer,
		client:           client,
		deleteClassifier: deleteClassifier,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string, samplesPath string) error {
	err := cmd.client.Create(classifierName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		fmt.Fprintln(cmd.writer, err.Error())
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	fmt.Fprintf(cmd.writer, "Adding samples from file '%s'... ", samplesPath)
	err = cmd.client.Train(classifierName, samplesPath)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		fmt.Fprintln(cmd.writer, err.Error())
		cmd.deleteClassifier.Execute(classifierName)
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}
