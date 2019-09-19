package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/waives/surf/config"
	"github.com/waives/surf/ioutils"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
)

//go:generate mockery -name "ClassifierUploader"
type ClassifierUploader interface {
	Upload(ctx context.Context, name string, contents io.Reader) error
}

func ConfigureUploadClassifierCommand(ctx context.Context,
	uploadCmd *kingpin.CmdClause,
	globalFlags *config.GlobalFlags) {
	args := &uploadClassifierArgs{}
	cmd := UploadClassifierCmd{}

	uploadClassifierCli := uploadCmd.Command("classifier",
		"Upload waives classifier (.clf file).").
		Action(func(parseContext *kingpin.ParseContext) error {
			msg := fmt.Sprintf("Creating classifier '%s' from '%s'... ", args.name, args.classifierFile)
			return ExecuteWithMessage(msg, func() error {
				err := cmd.initFromArgs(args, globalFlags)
				if err != nil {
					return err
				}

				return cmd.Execute(ctx)
			})
		})

	uploadClassifierCli.
		Arg("name", "The name of the new classifier.").
		Required().
		StringVar(&args.name)
	uploadClassifierCli.
		Arg("classifier-file", "The trained classifier file.").
		Required().
		StringVar(&args.classifierFile)
}

type uploadClassifierArgs struct {
	name           string
	classifierFile string
}

type UploadClassifierCmd struct {
	Uploader           ClassifierUploader
	ClassifierName     string
	ClassifierContents io.Reader
}

// Execute runs the command.
func (cmd *UploadClassifierCmd) Execute(ctx context.Context) error {
	defer ioutils.TryClose(cmd.ClassifierContents)

	return cmd.Uploader.Upload(ctx, cmd.ClassifierName, cmd.ClassifierContents)
}

func (cmd *UploadClassifierCmd) initFromArgs(args *uploadClassifierArgs, flags *config.GlobalFlags) error {
	var err error
	cmd.ClassifierContents, err = os.Open(args.classifierFile)
	if err != nil {
		return errors.New(fmt.Sprintf("the file '%s' could not be found", args.classifierFile))
	}

	cmd.ClassifierName = args.name

	apiClient, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Uploader = apiClient.Classifiers

	return nil
}
