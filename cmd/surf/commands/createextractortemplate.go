package commands

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/config"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"strings"
)

// CreateExtractorTemplate is a command which accepts a list of module ids,
// retrieves their descriptions from waives and generates a template json
// from them.
type CreateExtractorTemplateCmd struct {
	Client    ModuleGetter
	ModuleIds []string
	Output    io.Writer
}

type createExtractorTemplateArgs struct {
	moduleIds []string
}

// ConfigureCreateExtractorTemplateCmd configures kingpin with the 'create extractor-template'
// command.
func ConfigureCreateExtractorTemplateCmd(ctx context.Context, createCmd *kingpin.CmdClause,
	flags *config.GlobalFlags) {
	args := &createExtractorTemplateArgs{}
	createExtractorTemplateCmd := &CreateExtractorTemplateCmd{}

	createExtractorTemplateCli := createCmd.Command("extractor-template",
		"Create an extractor template from the provided module ids.").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := createExtractorTemplateCmd.initFromArgs(args, flags)

			if err != nil {
				return err
			}

			return createExtractorTemplateCmd.Execute(ctx)
		})

	createExtractorTemplateCli.
		Arg("module-ids", "The module IDs to include in the template").
		Required().
		StringsVar(&args.moduleIds)
}

// Execute runs the command.
func (cmd CreateExtractorTemplateCmd) Execute(ctx context.Context) error {
	var (
		jsonData []byte
		err      error
	)

	if len(cmd.ModuleIds) == 0 {
		return errors.New("at least one module ID must be specified")
	}

	allModules, err := cmd.Client.GetAll(ctx)
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
		return errors.WithMessage(err, "unable to create template")
	}

	_, err = cmd.Output.Write(jsonData)

	return err
}

func (cmd CreateExtractorTemplateCmd) getSpecifiedModules(existingModules ch360.ModuleList) (ch360.ModuleList, error) {
	var (
		missingModules []string
		presentModules ch360.ModuleList
	)

	// annoyingly we can't use a map here, since we want a case-insensitive search.
	for _, requestedModuleID := range cmd.ModuleIds {
		if existingModule := existingModules.Find(requestedModuleID); existingModule != nil {
			presentModules = append(presentModules, *existingModule)
		} else {
			missingModules = append(missingModules, requestedModuleID)
		}
	}

	if len(missingModules) > 0 {
		return nil, errors.Errorf("the following modules could not be found: %s",
			strings.Join(missingModules, ", "))
	}

	return presentModules, nil
}

// buildExtractorTemplateFor builds a ch360.ExtractorTemplate instance from a
// specified ch360.ModuleList.
func (cmd CreateExtractorTemplateCmd) buildExtractorTemplateFor(modules ch360.ModuleList) ch360.ExtractorTemplate {
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

func (cmd *CreateExtractorTemplateCmd) initFromArgs(args *createExtractorTemplateArgs,
	flags *config.GlobalFlags) error {
	cmd.ModuleIds = args.moduleIds

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = client.Modules
	cmd.Output = os.Stdout

	return nil
}
