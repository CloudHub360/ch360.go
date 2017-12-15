package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
)

//go:generate mockery -name "Deleter|Getter|DeleterGetter|ClassifierCommand"

type Deleter interface {
	Delete(name string) error
}

type Getter interface {
	GetAll() (ch360.ClassifierList, error)
}

type DeleterGetter interface {
	Deleter
	Getter
}

type ClassifierCommand interface {
	Execute(classifierName string) error
}

type DeleteClassifier struct {
	client DeleterGetter
}

func NewDeleteClassifier(client DeleterGetter) ClassifierCommand {
	return &DeleteClassifier{
		client: client,
	}
}

func (cmd *DeleteClassifier) Execute(classifierName string) error {
	classifiers, err := cmd.client.GetAll()

	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}

	if !classifiers.Contains(classifierName) {
		fmt.Println("[FAILED]")
		return errors.New("There is no classifier named '" + classifierName + "'")
	}
	return cmd.client.Delete(classifierName)
}
