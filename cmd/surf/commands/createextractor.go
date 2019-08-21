package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/net"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

type CreateExtractorCmd struct {
	Creator       ExtractorCreator
	ExtractorName string
	Template      *ch360.ExtractorTemplate
}

type createExtractorArgs struct {
	extractorName    string
	moduleIds        []string
	templateFilename string
}

// ConfigureCreateExtractorCmd configures kingpin with the 'create extractor' commands.
func ConfigureCreateExtractorCmd(ctx context.Context, createCmd *kingpin.CmdClause,
	flags *config.GlobalFlags) {
	args := &createExtractorArgs{}
	createExtractorCmd := &CreateExtractorCmd{}
	createExtractorCli := createCmd.Command("extractor", "Create waives extractor.")

	createExtractorFromModulesCli := createExtractorCli.
		Command("from-modules", "Create waives extractor from a set of modules.")
	createExtractorFromModulesCli.
		Arg("name", "The name of the new extractor.").
		Required().
		StringVar(&args.extractorName)
	createExtractorFromModulesCli.
		Arg("module-ids", "The module ids to create the extractor from.").
		Required().
		StringsVar(&args.moduleIds)

	createExtractorFromTemplateCli := createExtractorCli.Command("from-template",
		"The extractor template to create the extractor from.")
	createExtractorFromTemplateCli.
		Arg("name", "The name of the new extractor.").
		Required().
		StringVar(&args.extractorName)
	createExtractorFromTemplateCli.
		Arg("template-file", "The extraction template file (json).").
		Required().
		StringVar(&args.templateFilename)

	createExtractorFromModulesCli.
		Action(func(parseContext *kingpin.ParseContext) error {
			msg := fmt.Sprintf("Creating extractor '%s'... ", args.extractorName)

			return ExecuteWithMessage(msg, func() error {
				err := createExtractorCmd.initFromModuleIdArgs(args, flags)
				if err != nil {
					return err
				}
				return createExtractorCmd.Execute(ctx)
			})
		})

	createExtractorFromTemplateCli.
		Action(func(parseContext *kingpin.ParseContext) error {
			msg := fmt.Sprintf("Creating extractor '%s'... ", args.extractorName)
			return ExecuteWithMessage(msg, func() error {
				err := createExtractorCmd.initFromTemplateArgs(args, flags)
				if err != nil {

					return err
				}

				return createExtractorCmd.Execute(ctx)
			})
		})
}

// Execute runs the 'create extractor' command.
func (cmd *CreateExtractorCmd) Execute(ctx context.Context) error {
	err := cmd.Creator.CreateFromModules(ctx, cmd.ExtractorName, *cmd.Template)

	if err != nil {
		if detailedResponse, ok := err.(*net.DetailedErrorResponse); ok {
			return buildDetailedErrorMessage(*detailedResponse)
		}
	}

	return err
}

func (cmd *CreateExtractorCmd) initFromModuleIdArgs(args *createExtractorArgs, flags *config.GlobalFlags) error {
	var template = new(ch360.ExtractorTemplate)

	for _, moduleId := range args.moduleIds {
		template.Modules = append(template.Modules, ch360.ModuleTemplate{
			ID: moduleId,
		})
	}

	cmd.Template = template

	return cmd.initFromArgs(args, flags)
}

func (cmd *CreateExtractorCmd) initFromTemplateArgs(args *createExtractorArgs, flags *config.GlobalFlags) error {
	templateFile, err := os.Open(args.templateFilename)

	if err != nil {
		// err is guaranteed to be os.PathError
		pathErr := err.(*os.PathError)
		return errors.Errorf("Error when opening template file '%s': %v", args.templateFilename, pathErr.Err.Error())
	}

	cmd.Template, err = ch360.NewModulesTemplateFromJson(templateFile)

	if err != nil {
		return errors.WithMessagef(err, "Error when reading json template '%s'", args.templateFilename)
	}

	return cmd.initFromArgs(args, flags)
}

func (cmd *CreateExtractorCmd) initFromArgs(args *createExtractorArgs, flags *config.GlobalFlags) error {
	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Creator = client.Extractors
	cmd.ExtractorName = args.extractorName
	return nil
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
