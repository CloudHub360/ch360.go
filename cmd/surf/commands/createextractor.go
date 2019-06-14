package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/pkg/errors"
	"io"
	"os"
)

const CreateExtractorCommand = "create extractor"

type CreateExtractor struct {
	writer        io.Writer
	creator       ExtractorCreator
	extractorName string
	template      *ch360.ModulesTemplate
}

func NewCreateExtractor(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	template *ch360.ModulesTemplate) *CreateExtractor {
	return &CreateExtractor{
		writer:        writer,
		creator:       creator,
		template:      template,
		extractorName: extractorName,
	}
}

func NewCreateExtractorWithArgs(params *config.RunParams,
	client ExtractorCreator, out io.Writer) (*CreateExtractor, error) {

	var (
		err          error
		templateFile *os.File
		template     = new(ch360.ModulesTemplate)
	)

	if params.ModulesTemplate != "" {
		// deserialise template file
		templateFile, err = os.Open(params.ModulesTemplate)

		if err != nil {
			return nil, err
		}

		template, err = ch360.NewModulesTemplateFromJson(templateFile)

		if err != nil {
			return nil, errors.WithMessagef(err, "Error when reading json template '%s'", params.ModulesTemplate)
		}

	} else {
		// create template struct from IDs
		for _, moduleId := range params.ModuleIds {
			template.Modules = append(template.Modules, ch360.ModuleTemplate{
				ID: moduleId,
			})
		}
	}

	return NewCreateExtractor(out,
		client,
		params.Name,
		template), nil
}

func (cmd *CreateExtractor) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Creating extractor '%s'... ", cmd.extractorName)

	err := cmd.creator.CreateFromModules(ctx, cmd.extractorName, *cmd.template)

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
