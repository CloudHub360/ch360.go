package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"io"
)

type CreateExtractorFromModuleIds struct {
	writer        io.Writer
	creator       ExtractorCreator
	extractorName string
	moduleIds     []string
}

func NewCreateExtractorFromModuleIds(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	moduleIds []string) *CreateExtractorFromModuleIds {
	return &CreateExtractorFromModuleIds{
		writer:        writer,
		creator:       creator,
		moduleIds:     moduleIds,
		extractorName: extractorName,
	}
}

func NewCreateExtractorFromModuleIdsWithArgs(params *config.RunParams,
	client ExtractorCreator, out io.Writer) (*CreateExtractorFromModuleIds, error) {

	return NewCreateExtractorFromModuleIds(out,
		client,
		params.Name,
		params.ModuleIds), nil
}

func (cmd *CreateExtractorFromModuleIds) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Creating extractor '%s'... ", cmd.extractorName)

	template := ch360.ModulesTemplate{
		Modules: []ch360.ModuleTemplate{},
	}

	for _, moduleId := range cmd.moduleIds {
		template.Modules = append(template.Modules, ch360.ModuleTemplate{
			ID: moduleId,
		})
	}

	err := cmd.creator.CreateFromModules(ctx, cmd.extractorName, template)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}

func (cmd CreateExtractorFromModuleIds) Usage() string {
	return CreateExtractorCommand
}
