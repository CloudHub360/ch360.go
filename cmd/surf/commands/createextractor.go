package commands

import (
	"context"
	"fmt"
	"github.com/docopt/docopt-go"
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

func NewCreateExtractorFromArgs(args docopt.Opts,
	client ExtractorCreator, out io.Writer) (*CreateExtractor, error) {
	var (
		extractorName, _ = args.String("<name>")
		configPath, _    = args.String("<config-file>")
	)

	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	return NewCreateExtractor(out,
		client,
		extractorName,
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
