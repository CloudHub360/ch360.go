package commands

import (
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/pkg/errors"
	"io"
	"os"
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
	return ExecuteWithMessage("Creating extractor template...", os.Stderr, func() error {
		defer cmd.close()

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

		jsonData, err := json.MarshalIndent(template, "", "  ")

		if err != nil {
			return errors.WithMessage(err, "Unable to create template")
		}

		_, err = cmd.writer.Write(jsonData)

		return err
	})
}

// Attempts to close any opened files
func (cmd CreateExtractorTemplate) close() {
	if closer, ok := cmd.writer.(io.WriteCloser); ok {
		_ = closer.Close()
	}
}

func (cmd CreateExtractorTemplate) getSpecifiedModules(existingModules ch360.ModuleList) (ch360.ModuleList, error) {
	modulesMap := existingModules.Map()

	var (
		missingModules = []string{}
		presentModules ch360.ModuleList
	)

	for _, requestedModuleID := range cmd.moduleIds {
		if presentModule, ok := modulesMap[requestedModuleID]; ok {
			presentModules = append(presentModules, presentModule)
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

func NewCreateExtractorTemplate(moduleIds []string, client ModuleGetter, out io.Writer) (*CreateExtractorTemplate, error) {
	return &CreateExtractorTemplate{
		moduleIds: moduleIds,
		client:    client,
		writer:    out,
	}, nil
}

func NewCreateExtractorTemplateWithArgs(params *config.RunParams, client ModuleGetter) (*CreateExtractorTemplate, error) {

	var (
		out = os.Stdout
		err error
	)

	if params.OutputFile != "" {
		out, err = os.Create(params.OutputFile)
	}

	if err != nil {
		// os.Create's err is guaranteed to be os.PathError, so we
		// get the underlying cause from it to avoid duplicating the path
		// in the message output to the user
		err := err.(*os.PathError).Err
		return nil, errors.WithMessagef(err, "Could not open file '%s'", params.OutputFile)
	}

	return NewCreateExtractorTemplate(params.ModuleIds, client, out)
}
