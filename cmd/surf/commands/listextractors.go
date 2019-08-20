package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

//go:generate mockery -name "ExtractorDeleter|ExtractorGetter|ExtractorCommand"

const ListExtractorsCommand = "list extractors"

type ExtractorDeleter interface {
	Delete(ctx context.Context, name string) error
}

type ExtractorGetter interface {
	GetAll(ctx context.Context) (ch360.ExtractorList, error)
}

type ListExtractorsCmd struct {
	Client ExtractorGetter
}

// ConfigureListExtractorsCmd configures kingpin with the 'list extractors' command.
func ConfigureListExtractorsCmd(ctx context.Context, listCmd *kingpin.CmdClause, flags *config.GlobalFlags) {
	listExtractorsCmd := &ListExtractorsCmd{}

	listCmd.Command("extractors", "List all available extractors.").
		Action(func(parseContext *kingpin.ParseContext) error {
			exitOnErr(listExtractorsCmd.initFromArgs(flags))
			exitOnErr(listExtractorsCmd.Execute(ctx))
			return nil
		})
}

// Execute runs the 'list extractors' command.
func (cmd *ListExtractorsCmd) Execute(ctx context.Context) error {
	extractors, err := cmd.Client.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(extractors) == 0 {
		fmt.Println("No extractors found.")
	}

	for _, extractor := range extractors {
		fmt.Println(extractor.Name)
	}

	return nil
}

func (cmd ListExtractorsCmd) Usage() string {
	return ListExtractorsCommand
}

func (cmd *ListExtractorsCmd) initFromArgs(flags *config.GlobalFlags) error {
	apiClient, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = apiClient.Extractors
	return nil
}
