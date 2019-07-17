package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"io"
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

func NewCreateExtractorFromModules(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	modules []string) *CreateExtractor {

	var template = new(ch360.ExtractorTemplate)

	for _, moduleId := range modules {
		template.Modules = append(template.Modules, ch360.ModuleTemplate{
			ID: moduleId,
		})
	}

	return NewCreateExtractor(writer,
		creator,
		extractorName,
		template)
}

func NewCreateExtractorFromTemplate(out io.Writer,
	client ExtractorCreator, extractorName string, templateFile io.Reader) (*CreateExtractor,
	error) {

	var (
		err      error
		template = new(ch360.ExtractorTemplate)
	)

	template, err = ch360.NewModulesTemplateFromJson(templateFile)

	if err != nil {
		return nil, errors.WithMessagef(err,
			"Error when reading json template")
	}

	return NewCreateExtractor(out,
		client,
		extractorName,
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
