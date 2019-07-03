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
	template      *ch360.ExtractorTemplate
}

func NewCreateExtractor(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	template *ch360.ExtractorTemplate) *CreateExtractor {
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
		template     = new(ch360.ExtractorTemplate)
	)

	if params.ModulesTemplate != "" {
		// deserialise template file
		templateFile, err = os.Open(params.ModulesTemplate)

		if err != nil {
			pathErr := err.(*os.PathError)
			return nil, errors.WithMessagef(pathErr.Err,
				"Error when opening template file '%s'", params.ModulesTemplate)
		}

		template, err = ch360.NewModulesTemplateFromJson(templateFile)

		if err != nil {
			return nil, errors.WithMessagef(err,
				"Error when reading json template '%s'", params.ModulesTemplate)
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
	return ExecuteWithMessage(fmt.Sprintf("Creating extractor '%s'... ", cmd.extractorName), func() error {
		return cmd.creator.CreateFromModules(ctx, cmd.extractorName, *cmd.template)
	})

}

func (cmd CreateExtractor) Usage() string {
	return CreateExtractorCommand
}
