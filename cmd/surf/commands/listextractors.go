package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
)

//go:generate mockery -name "ExtractorDeleter|ExtractorGetter|ExtractorDeleterGetter|ExtractorCommand"

const ListExtractorsCommand = "list extractors"

type ExtractorDeleter interface {
	Delete(ctx context.Context, name string) error
}

type ExtractorGetter interface {
	GetAll(ctx context.Context) (ch360.ExtractorList, error)
}

type ExtractorDeleterGetter interface {
	ExtractorDeleter
	ExtractorGetter
}

type ListExtractors struct {
	client ExtractorGetter
	writer io.Writer
}

func NewListExtractors(client ExtractorGetter, out io.Writer) *ListExtractors {
	return &ListExtractors{
		client: client,
		writer: out,
	}
}

func (cmd *ListExtractors) Execute(ctx context.Context) error {
	extractors, err := cmd.client.GetAll(ctx)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if len(extractors) == 0 {
		fmt.Fprintln(cmd.writer, "No extractors found.")
	}

	for _, extractor := range extractors {
		fmt.Fprintln(cmd.writer, extractor.Name)
	}

	return nil
}

func (cmd ListExtractors) Usage() string {
	return ListExtractorsCommand
}
