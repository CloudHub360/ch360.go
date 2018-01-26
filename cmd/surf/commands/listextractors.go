package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
)

//go:generate mockery -name "ExtractorDeleter|ExtractorGetter|ExtractorDeleterGetter|ExtractorCommand"

type ExtractorDeleter interface {
	Delete(name string) error
}

type ExtractorGetter interface {
	GetAll() (ch360.ExtractorList, error)
}

type ExtractorDeleterGetter interface {
	ExtractorDeleter
	ExtractorGetter
}

type ListExtractors struct {
	client ExtractorGetter
	writer io.Writer
}

func NewListExtractors(writer io.Writer, client ExtractorGetter) *ListExtractors {
	return &ListExtractors{
		client: client,
		writer: writer,
	}
}

func (cmd *ListExtractors) Execute() (ch360.ExtractorList, error) {
	extractors, err := cmd.client.GetAll()
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return nil, err
	}

	if len(extractors) == 0 {
		fmt.Fprintln(cmd.writer, "No extractors found.")
	}

	for _, extractor := range extractors {
		fmt.Fprintln(cmd.writer, extractor.Name)
	}

	return extractors, nil
}
