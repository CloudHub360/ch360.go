package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
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
		err := cmd.creator.CreateFromModules(ctx, cmd.extractorName, *cmd.template)

		if err != nil {
			if detailedResponse, ok := err.(*net.DetailedErrorResponse); ok {
				return buildDetailedErrorMessage(*detailedResponse)
			}
		}

		return err
	})

}

func (cmd CreateExtractor) Usage() string {
	return CreateExtractorCommand
}

func buildDetailedErrorMessage(errorResponse net.DetailedErrorResponse) error {
	//noinspection ALL odd names to match json
	type detailedError struct {
		Module_ID      string
		Messages       []string
		Path           string
		Argument_Name  string
		Argument_Value string
	}

	var detailedErrs []detailedError
	err := mapstructure.Decode(errorResponse.Errors, &detailedErrs)

	if err != nil {
		return errors.WithMessage(&errorResponse, "Could not deserialise response from server")
	}

	sb := strings.Builder{}
	sb.WriteString("Extractor creation failed with the following error: ")
	sb.WriteString(fmt.Sprintf("%s\n", errorResponse.Error()))

	// group error info by module
	errorsByModule := map[string][]detailedError{}
	for _, detailedErr := range detailedErrs {
		moduleId := detailedErr.Module_ID
		errorsByModule[moduleId] = append(errorsByModule[moduleId], detailedErr)
	}

	for moduleId, detailedErrs := range errorsByModule {
		if moduleId == "" {
			moduleId = "(not found)"
		}

		sb.WriteString(fmt.Sprintf("\nModule %s:\n", moduleId))
		for _, detailedErr := range detailedErrs {

			if detailedErr.Argument_Name != "" {
				// param err
				for _, message := range detailedErr.Messages {
					sb.WriteString(fmt.Sprintf("  Parameter \"%s\": %s (specified \"%s\")\n",
						detailedErr.Argument_Name,
						message,
						detailedErr.Argument_Value))
				}
			} else {
				// module err
				sb.WriteString(fmt.Sprintf("  %s\n", strings.Join(detailedErr.Messages, ", ")))
			}
		}
	}

	return errors.New(sb.String())
}
