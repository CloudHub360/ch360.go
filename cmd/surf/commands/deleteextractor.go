package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"io"
)

const DeleteExtractorCommand = "delete extractor"

type DeleteExtractor struct {
	client        ExtractorDeleterGetter
	writer        io.Writer
	extractorName string
}

func NewDeleteExtractor(extractorName string, writer io.Writer, client ExtractorDeleterGetter) *DeleteExtractor {
	return &DeleteExtractor{
		writer:        writer,
		client:        client,
		extractorName: extractorName,
	}
}

func NewDeleteExtractorFromArgs(args docopt.Opts, client ExtractorDeleterGetter, writer io.Writer) (*DeleteExtractor, error) {
	extractorName, err := args.String("<name>")

	if err != nil {
		return nil, err
	}

	return NewDeleteExtractor(extractorName, writer, client), nil
}

func (cmd *DeleteExtractor) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Deleting extractor '%s'... ", cmd.extractorName)

	extractors, err := cmd.client.GetAll()

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if !extractors.Contains(cmd.extractorName) {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return errors.New("There is no extractor named '" + cmd.extractorName + "'")
	}

	err = cmd.client.Delete(cmd.extractorName)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")
	return nil
}

func (cmd *DeleteExtractor) Usage() string {
	return DeleteExtractorCommand
}
