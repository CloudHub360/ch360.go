package commands

import (
	"fmt"
	"os"
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
	client           CreatorTrainer
	deleteClassifier ClassifierCommand
}

func NewCreateClassifier(client CreatorTrainer, deleteClassifier ClassifierCommand) *CreateClassifier {
	return &CreateClassifier{
		client:           client,
		deleteClassifier: deleteClassifier,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string, samplesPath string) error {
	err := cmd.client.Create(classifierName)
	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}

	fmt.Println("[OK]")

	fmt.Printf("Adding samples from file '%s'... ", samplesPath)
	err = cmd.client.Train(classifierName, samplesPath)

	if err != nil {
		fmt.Println("[FAILED]")
		fmt.Fprint(os.Stderr, err.Error())
		cmd.deleteClassifier.Execute(classifierName)
		return err
	}

	fmt.Println("[OK]")
	return nil
}
