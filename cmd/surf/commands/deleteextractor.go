package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/waives/surf/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

//go:generate mockery -name "ExtractorDeleterGetter"
type ExtractorDeleterGetter interface {
	ExtractorDeleter
	ExtractorGetter
}

type deleteExtractorArgs struct {
	extractorName string
}

type DeleteExtractorCmd struct {
	Client        ExtractorDeleterGetter
	ExtractorName string
}

// ConfigureDeleteExtractorCmd configures kingpin with the 'delete extractor' command.
func ConfigureDeleteExtractorCmd(ctx context.Context, deleteCmd *kingpin.CmdClause, flags *config.
	GlobalFlags) {
	args := &deleteExtractorArgs{}
	deleteExtractorCmd := &DeleteExtractorCmd{}

	deleteExtractorCli := deleteCmd.Command("extractor", "Delete waives extractor.").
		Action(func(parseContext *kingpin.ParseContext) error {
			return ExecuteWithMessage(fmt.Sprintf("Deleting extractor '%s'... ", args.extractorName),
				func() error {

					err := deleteExtractorCmd.initFromArgs(args, flags)
					if err != nil {

						return err
					}
					return deleteExtractorCmd.Execute(ctx)
				})
		})

	deleteExtractorCli.
		Arg("name", "The name of the extractor to delete.").
		Required().
		StringVar(&args.extractorName)
}

func (cmd *DeleteExtractorCmd) Execute(ctx context.Context) error {
	extractors, err := cmd.Client.GetAll(ctx)

	if err != nil {
		return err
	}

	if !extractors.Contains(cmd.ExtractorName) {
		return errors.New("there is no extractor named '" + cmd.ExtractorName + "'")
	}

	return cmd.Client.Delete(ctx, cmd.ExtractorName)
}

func (cmd *DeleteExtractorCmd) initFromArgs(args *deleteExtractorArgs, flags *config.GlobalFlags) error {
	cmd.ExtractorName = args.extractorName

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = client.Extractors
	return nil
}
