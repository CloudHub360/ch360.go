package commands

import (
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
)

//go:generate mockery -name "Deleter|Getter|DeleteGetter"

func NewCreateClassifier(client *ch360.ClassifiersClient) *CreateClassifier {
	return &CreateClassifier{
		client: client,
	}
type Deleter interface {
	Delete (name string) error
}

type Getter interface {
	GetAll () ([]ch360.Classifier, error)
}

type DeleteGetter interface {
	Deleter
	Getter
}

type DeleteClassifier struct {
	client DeleteGetter
}

func NewDeleteClassifier(client DeleteGetter) *DeleteClassifier {
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
	return cmd.client.Delete(classifierName)
}
