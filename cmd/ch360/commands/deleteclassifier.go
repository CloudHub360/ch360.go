package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
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
	writer io.Writer
}

func NewDeleteClassifier(writer io.Writer, client DeleterGetter) ClassifierCommand {
	return &DeleteClassifier{
		writer: writer,
		client: client,
	}
}

func (cmd *DeleteClassifier) Execute(classifierName string) error {
	classifiers, err := cmd.client.GetAll()

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if !classifiers.Contains(classifierName) {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return errors.New("There is no classifier named '" + classifierName + "'")
	}

	err = cmd.client.Delete(classifierName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	return nil
}
