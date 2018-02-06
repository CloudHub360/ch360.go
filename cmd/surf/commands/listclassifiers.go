package commands

import (
	"context"
	"fmt"
	"io"
)

const ListClassifiersCommand = "list classifiers"

type ListClassifiers struct {
	client ClassifierGetter
	writer io.Writer
}

func NewListClassifiers(client ClassifierGetter, out io.Writer) *ListClassifiers {
	return &ListClassifiers{
		client: client,
		writer: out,
	}
}

func (cmd *ListClassifiers) Execute(ctx context.Context) error {
	classifiers, err := cmd.client.GetAll()
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if !classifiers.Any() {
		fmt.Fprintln(cmd.writer, "No classifiers found.")
	}

	for _, classifier := range classifiers {
		fmt.Fprintln(cmd.writer, classifier.Name)
	}

	return nil
}
