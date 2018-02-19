package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"io"
	"os"
)

//go:generate mockery -name "ExtractorCreator"

const CreateExtractorCommand = "create extractor"

type ExtractorCreator interface {
	Create(name string, config io.Reader) error
}

type CreateExtractor struct {
	writer        io.Writer
	creator       ExtractorCreator
	extractorName string
	config        io.Reader
}

func NewCreateExtractor(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	config io.Reader) *CreateExtractor {
	return &CreateExtractor{
		writer:        writer,
		creator:       creator,
		config:        config,
		extractorName: extractorName,
	}
}

func NewCreateExtractorFromArgs(params *config.RunParams,
	client ExtractorCreator, out io.Writer) (*CreateExtractor, error) {

	configFile, err := os.Open(params.ConfigPath)

	if err != nil {
		return nil, err
	}

	return NewCreateExtractor(out,
		client,
		params.ExtractorName,
		configFile), nil
}

func (cmd *CreateExtractor) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Creating extractor '%s'... ", cmd.extractorName)

	err := cmd.creator.Create(cmd.extractorName, cmd.config)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}

func (cmd CreateExtractor) Usage() string {
	return CreateExtractorCommand
}
