package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
)

type ListClassifiers struct {
	client ClassifierGetter
	writer io.Writer
}

func NewListClassifiers(writer io.Writer, client ClassifierGetter) *ListClassifiers {
	return &ListClassifiers{
		client: client,
		writer: writer,
	}
}

func (cmd *ListClassifiers) Execute() (ch360.ClassifierList, error) {
	classifiers, err := cmd.client.GetAll()
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return nil, err
	}

	if !classifiers.Any() {
		fmt.Fprintln(cmd.writer, "No classifiers found.")
	}

	for _, classifier := range classifiers {
		fmt.Fprintln(cmd.writer, classifier.Name)
	}

	return classifiers, nil
}
