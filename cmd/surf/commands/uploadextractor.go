package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/ioutils"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
)

//go:generate mockery -name "ExtractorCreator"

type ExtractorCreator interface {
	Create(ctx context.Context, name string, config io.Reader) error
	CreateFromJson(ctx context.Context, name string, jsonTemplate io.Reader) error
	CreateFromModules(ctx context.Context, name string, modules ch360.ExtractorTemplate) error
}

type UploadExtractorCmd struct {
	ExtractorCreator ExtractorCreator
	ExtractorName    string
	ExtractorContent io.Reader
}

type uploadExtractorArgs struct {
	extractorName string
	extractorFile string
}

func ConfigureUploadExtractorCommand(ctx context.Context,
	uploadCommand *kingpin.CmdClause,
	globalFlags *config.GlobalFlags) {
	args := &uploadExtractorArgs{}

	cmd := &UploadExtractorCmd{}
	uploadExtractorCommand := uploadCommand.
		Command("extractor", "Upload waives extractor (.fpxlc file).")
	uploadExtractorCommand.
		Arg("name", "The name of the new extractor.").
		Required().
		StringVar(&args.extractorName)
	uploadExtractorCommand.
		Arg("config-file", "The extraction configuration file.").
		Required().
		StringVar(&args.extractorFile)

	uploadExtractorCommand.Action(func(parseContext *kingpin.ParseContext) error {
		return ExecuteWithMessage(
			fmt.Sprintf("Uploading extractor '%s'... ", args.extractorName),
			func() error {
				err := cmd.initFromArgs(args, globalFlags)
				if err != nil {
					return err
				}

				return cmd.Execute(ctx)
			})
	})
}

func (cmd *UploadExtractorCmd) Execute(ctx context.Context) error {
	defer ioutils.TryClose(cmd.ExtractorContent)

	return cmd.ExtractorCreator.Create(ctx, cmd.ExtractorName, cmd.ExtractorContent)
}

func (cmd *UploadExtractorCmd) initFromArgs(args *uploadExtractorArgs,
	flags *config.GlobalFlags) error {

	var err error
	cmd.ExtractorContent, err = os.Open(args.extractorFile)
	if err != nil {
		return errors.Errorf("the file '%s' could not be found", args.extractorFile)
	}

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.ExtractorCreator = client.Extractors
	cmd.ExtractorName = args.extractorName
	return nil
}
