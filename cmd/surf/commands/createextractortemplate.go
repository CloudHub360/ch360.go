package commands

import (
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ioutils"
	"github.com/pkg/errors"
	"io"
	"strings"
)

const CreateExtractorTemplateCommand = "create extractor-template"

var _ Command = (*CreateExtractorTemplate)(nil)

// CreateExtractorTemplate is a command which accepts a list of module ids,
// retrieves their descriptions from waives and generates a template json
// from them.
type CreateExtractorTemplate struct {
	client    ModuleGetter
	writer    io.Writer
	moduleIds []string
}

// Execute runs the command.
func (cmd CreateExtractorTemplate) Execute(ctx context.Context) error {
	defer ioutils.TryClose(cmd.writer)
	var (
		jsonData []byte
		err      error
	)

	err = ExecuteWithMessage("Creating extractor template...", func() error {
		if len(cmd.moduleIds) == 0 {
			return errors.New("At least one module ID must be specified")
		}

		allModules, err := cmd.client.GetAll(ctx)
		if err != nil {
			return err
		}

		specifedModules, err := cmd.getSpecifiedModules(allModules)
		if err != nil {
			return err
		}

		template := cmd.buildExtractorTemplateFor(specifedModules)

		jsonData, err = json.MarshalIndent(template, "", "  ")

		if err != nil {
			return errors.WithMessage(err, "Unable to create template")
		}

		return nil
	})

	if err != nil {
		return err
	}

	_, err = cmd.writer.Write(jsonData)

	return err
}

func (cmd CreateExtractorTemplate) getSpecifiedModules(existingModules ch360.ModuleList) (ch360.ModuleList, error) {
	var (
		missingModules = []string{}
		presentModules ch360.ModuleList
	)

	// annoyingly we can't use a map here, since we want a case-insensitive search.
	for _, requestedModuleID := range cmd.moduleIds {
		if existingModule := existingModules.Find(requestedModuleID); existingModule != nil {
			presentModules = append(presentModules, *existingModule)
		} else {
			missingModules = append(missingModules, requestedModuleID)
		}
	}

	if len(missingModules) > 0 {
		return nil, errors.Errorf("The following modules could not be found: %s", strings.Join(missingModules, ", "))
	}

	return presentModules, nil
}

func (cmd CreateExtractorTemplate) Usage() string {
	return CreateExtractorTemplateCommand
}

// buildExtractorTemplateFor builds a ch360.ExtractorTemplate instance from a
// specified ch360.ModuleList.
func (cmd CreateExtractorTemplate) buildExtractorTemplateFor(modules ch360.ModuleList) ch360.ExtractorTemplate {
	template := ch360.ExtractorTemplate{}

	for _, module := range modules {
		argsMap := map[string]interface{}{}

		for _, arg := range module.Parameters {
			if !arg.Required {
				continue
			}

			argsMap[arg.ID] = ""
		}

		template.Modules = append(template.Modules, ch360.ModuleTemplate{
			ID:        module.ID,
			Arguments: argsMap,
		})
	}

	return template
}

func NewCreateExtractorTemplate(moduleIds []string, client ModuleGetter, out io.Writer) *CreateExtractorTemplate {
	return &CreateExtractorTemplate{
		moduleIds: moduleIds,
		client:    client,
		writer:    out,
	}
}
