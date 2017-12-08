package commands

import (
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
)

type CreateClassifier struct {
	client *ch360.ClassifiersClient
}

func NewCreateClassifier(client *ch360.ClassifiersClient) *CreateClassifier {
	return &CreateClassifier{
		client: client,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string) error {
	return cmd.client.CreateClassifier(classifierName)
}

type DeleteClassifier struct {
	client *ch360.ClassifiersClient
}

func NewDeleteClassifier(client *ch360.ClassifiersClient) *DeleteClassifier {
	return &DeleteClassifier{
		client: client,
	}
}

func (cmd *DeleteClassifier) Execute(classifierName string) error {
	classifiers, err := cmd.client.GetAll()

	if err != nil {
		return err
	}

	if !classifiers.Contains(classifierName) {
		return errors.New("There is no classifier named '" + classifierName + "'")
	}
	return cmd.client.DeleteClassifier(classifierName)
}
