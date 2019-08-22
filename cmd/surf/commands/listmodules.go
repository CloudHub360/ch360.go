package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

//go:generate mockery -name "ModuleGetter"

const ListModulesCommandString = "list modules"

type ModuleGetter interface {
	GetAll(ctx context.Context) (ch360.ModuleList, error)
}

type ListModulesCmd struct {
	Client ModuleGetter
}

func ConfigureListModulesCommand(ctx context.Context,
	listCmd *kingpin.CmdClause, globalFlags *config.GlobalFlags) {
	cmd := &ListModulesCmd{}

	listCmd.Command("modules", "List all available extractor modules.").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := cmd.initFromArgs(globalFlags)

			if err != nil {
				return err
			}
			return cmd.Execute(ctx)
		})
}

func (cmd *ListModulesCmd) initFromArgs(flags *config.GlobalFlags) error {
	var err error
	apiClient, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = apiClient.Modules
	return nil
}

func (cmd *ListModulesCmd) Execute(ctx context.Context) error {
	modules, err := cmd.Client.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		fmt.Println("No modules found.")
		return nil
	}

	table := NewTable(os.Stdout, []string{"Name", "ID", "Summary"})
	for _, module := range modules {
		table.Append([]string{module.Name, module.ID, module.Summary})
	}
	table.Render()

	return nil
}

func (cmd ListModulesCmd) Usage() string {
	return ListModulesCommandString
}
