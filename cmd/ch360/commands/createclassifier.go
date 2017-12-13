package commands

import "fmt"

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
	client CreatorTrainer
}

func NewCreateClassifier(client CreatorTrainer) *CreateClassifier {
	return &CreateClassifier{
		client: client,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string, samplesPath string) error {
	err := cmd.client.Create(classifierName)
	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}

	fmt.Println("[OK]")
	return cmd.client.Train(classifierName, samplesPath)
}
