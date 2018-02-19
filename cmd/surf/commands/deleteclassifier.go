package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"io"
)

//go:generate mockery -name "ClassifierDeleter|ClassifierGetter|ClassifierDeleterGetter|ClassifierCommand"

const DeleteClassifierCommand = "delete classifier"

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

type DeleteClassifier struct {
	client         ClassifierDeleterGetter
	writer         io.Writer
	classifierName string
}

func NewDeleteClassifier(classifierName string, writer io.Writer, client ClassifierDeleterGetter) *DeleteClassifier {
	return &DeleteClassifier{
		writer:         writer,
		client:         client,
		classifierName: classifierName,
	}
}

func NewDeleteClassifierFromArgs(params *config.RunParams, client ClassifierDeleterGetter, out io.Writer) (*DeleteClassifier, error) {
	return &DeleteClassifier{
		client:         client,
		writer:         out,
		classifierName: params.Name,
	}, nil
}

func (cmd *DeleteClassifier) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Deleting classifier '%s'... ", cmd.classifierName)

	classifiers, err := cmd.client.GetAll()

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if !classifiers.Contains(cmd.classifierName) {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return errors.New("There is no classifier named '" + cmd.classifierName + "'")
	}

	err = cmd.client.Delete(cmd.classifierName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}

func (cmd DeleteClassifier) Usage() string {
	return DeleteClassifierCommand
}
