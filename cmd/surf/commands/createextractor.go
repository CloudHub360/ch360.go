package commands

import (
	"fmt"
	"io"
)

//go:generate mockery -name "ExtractorCreator"

type ExtractorCreator interface {
	Create(name string, config io.Reader) error
}

type CreateExtractor struct {
	writer  io.Writer
	creator ExtractorCreator
}

func NewCreateExtractor(writer io.Writer,
	creator ExtractorCreator) *CreateExtractor {
	return &CreateExtractor{
		writer:  writer,
		creator: creator,
	}
}

func (cmd *CreateExtractor) Execute(extractorName string, config io.Reader) error {
	fmt.Fprintf(cmd.writer, "Creating extractor '%s'... ", extractorName)

	err := cmd.creator.Create(extractorName, config)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}
