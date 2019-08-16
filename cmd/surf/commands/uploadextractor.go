package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
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

type uploadExtractorArgs struct {
	extractorName string
	extractorFile *os.File
}

func ConfigureUploadExtractorCommand(ctx context.Context,
	uploadCommand *kingpin.CmdClause,
	globalFlags *config.GlobalFlags) {
	args := &uploadExtractorArgs{}

	uploadExtractorCommand := uploadCommand.
		Command("extractor", "Upload waives extractor (.fpxlc file).")
	uploadExtractorCommand.
		Arg("name", "The name of the new extractor.").
		Required().
		StringVar(&args.extractorName)
	uploadExtractorCommand.
		Arg("config-file", "The extraction configuration file.").
		Required().
		FileVar(&args.extractorFile)

	uploadExtractorCommand.Action(func(parseContext *kingpin.ParseContext) error {
		// execute the command
		return executeUploadExtractor(ctx, args, globalFlags)
	})
}

func executeUploadExtractor(ctx context.Context,
	args *uploadExtractorArgs,
	globalFlags *config.GlobalFlags) error {

	return ExecuteWithMessage(fmt.Sprintf("Uploading extractor '%s'... ", args.extractorName),
		func() error {
			client, err := initApiClient(globalFlags.ClientId,
				globalFlags.ClientSecret,
				globalFlags.LogHttp)

			if err != nil {
				return err
			}

			return client.Extractors.Create(ctx, args.extractorName, args.extractorFile)
		})
}
