package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"io"
	"os"
)

type CreateExtractorFromModules struct {
	writer        io.Writer
	creator       ExtractorCreator
	extractorName string
	template      io.ReadCloser
}

func NewCreateExtractorFromModules(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	templateFile io.ReadCloser) *CreateExtractorFromModules {
	return &CreateExtractorFromModules{
		writer:        writer,
		creator:       creator,
		template:      templateFile,
		extractorName: extractorName,
	}
}

func NewCreateExtractorFromModulesWithArgs(params *config.RunParams,
	client ExtractorCreator, out io.Writer) (*CreateExtractorFromModules, error) {

	configFile, err := os.Open(params.ModulesTemplate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be found.", params.ModulesTemplate))
	}

	return NewCreateExtractorFromModules(out,
		client,
		params.Name,
		configFile), nil
}

func (cmd *CreateExtractorFromModules) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Creating extractor '%s'... ", cmd.extractorName)

	err := cmd.creator.CreateFromModules(ctx, cmd.extractorName, cmd.template)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}

func (cmd CreateExtractorFromModules) Usage() string {
	return CreateExtractorCommand
}
