package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
)

//go:generate mockery -name "ClassifierDeleter|ClassifierGetter|ClassifierDeleterGetter|ClassifierCommand"

type ClassifierDeleter interface {
	Delete(name string) error
}

type ClassifierGetter interface {
	GetAll() (ch360.ClassifierList, error)
}

type ClassifierDeleterGetter interface {
	ClassifierDeleter
	ClassifierGetter
}

type ClassifierCommand interface {
	Execute(classifierName string) error
}

type DeleteClassifier struct {
	client ClassifierDeleterGetter
	writer io.Writer
}

func NewDeleteClassifier(writer io.Writer, client ClassifierDeleterGetter) ClassifierCommand {
	return &DeleteClassifier{
		writer: writer,
		client: client,
	}
}

func (cmd *DeleteClassifier) Execute(classifierName string) error {
	fmt.Fprintf(cmd.writer, "Deleting classifier '%s'... ", classifierName)

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

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}
