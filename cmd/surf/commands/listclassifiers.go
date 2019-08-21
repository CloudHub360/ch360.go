package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ListClassifiersCmd struct {
	Client ClassifierGetter
}

// Configures kingpin with the 'list classifiers' command
func ConfigureListClassifiersCmd(ctx context.Context, parentCmd *kingpin.CmdClause, flags *config.GlobalFlags) {
	listClassifiersCmd := &ListClassifiersCmd{}
	parentCmd.Command("classifiers", "List all available classifiers.").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := listClassifiersCmd.initFromArgs(flags)
			if err != nil {
				return err
			}
			return listClassifiersCmd.Execute(ctx)
		})
}

// Executes the command.
func (cmd *ListClassifiersCmd) Execute(ctx context.Context) error {
	classifiers, err := cmd.Client.GetAll(ctx)
	if err != nil {
		return err
	}

	if !classifiers.Any() {
		fmt.Println("No classifiers found.")
	}

	for _, classifier := range classifiers {
		fmt.Println(classifier.Name)
	}

	return nil
}

func (cmd *ListClassifiersCmd) initFromArgs(flags *config.GlobalFlags) error {
	var err error
	apiClient, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = apiClient.Classifiers
	return nil
}
