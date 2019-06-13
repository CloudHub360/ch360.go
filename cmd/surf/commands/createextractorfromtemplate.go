package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"io"
	"os"
)

type CreateExtractorFromTemplate struct {
	writer        io.Writer
	creator       ExtractorCreator
	extractorName string
	template      io.ReadCloser
}

func NewCreateExtractorFromTemplate(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	templateFile io.ReadCloser) *CreateExtractorFromTemplate {
	return &CreateExtractorFromTemplate{
		writer:        writer,
		creator:       creator,
		template:      templateFile,
		extractorName: extractorName,
	}
}

func NewCreateExtractorFromTemplateWithArgs(params *config.RunParams,
	client ExtractorCreator, out io.Writer) (*CreateExtractorFromTemplate, error) {

	configFile, err := os.Open(params.ModulesTemplate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be found.", params.ModulesTemplate))
	}

	return NewCreateExtractorFromTemplate(out,
		client,
		params.Name,
		configFile), nil
}

func (cmd *CreateExtractorFromTemplate) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Creating extractor '%s'... ", cmd.extractorName)

	err := cmd.creator.CreateFromJson(ctx, cmd.extractorName, cmd.template)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}

func (cmd CreateExtractorFromTemplate) Usage() string {
	return CreateExtractorCommand
}
