package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
)

//go:generate mockery -name "ClassifierCreator|ClassifierTrainer|ClassifierClient"

type ClassifierCreator interface {
	Create(ctx context.Context, name string) error
}

type ClassifierTrainer interface {
	Train(ctx context.Context, name string, samplesArchive io.Reader) error
}

type CreateClassifierCmd struct {
	Creator        ClassifierCreator
	Deleter        ClassifierDeleter
	Trainer        ClassifierTrainer
	ClassifierName string
	SamplesArchive *os.File
}

type createClassifierArgs struct {
	classifierName         string
	samplesArchiveFilename string
}

func ConfigureCreateClassifierCmd(ctx context.Context, createCmd *kingpin.CmdClause,
	flags *config.GlobalFlags) {
	args := &createClassifierArgs{}
	createClassifierCmd := &CreateClassifierCmd{}

	createClassifierCli := createCmd.Command("classifier", "Create waives classifier from a set of samples.").
		Action(func(parseContext *kingpin.ParseContext) error {
			return ExecuteWithMessage(fmt.Sprintf("Creating classifier '%s'... ",
				args.classifierName),
				func() error {
					err := createClassifierCmd.initFromArgs(args, flags)

					if err != nil {
						return err
					}

					return createClassifierCmd.Execute(ctx)
				})
		})

	createClassifierCli.
		Arg("name", "The name of the new classifier.").
		Required().
		StringVar(&args.classifierName)

	createClassifierCli.
		Arg("samples-zip", "The zip file containing training samples.").
		Required().
		StringVar(&args.samplesArchiveFilename)
}

func (cmd *CreateClassifierCmd) Execute(ctx context.Context) error {
	defer cmd.SamplesArchive.Close()

	err := cmd.Creator.Create(ctx, cmd.ClassifierName)
	if err != nil {
		return err
	}

	err = cmd.Trainer.Train(ctx, cmd.ClassifierName, cmd.SamplesArchive)

	if err != nil {
		_ = cmd.Deleter.Delete(ctx, cmd.ClassifierName)
		return err
	}

	return nil
}

func (cmd *CreateClassifierCmd) initFromArgs(args *createClassifierArgs, flags *config.GlobalFlags) error {
	var err error
	cmd.SamplesArchive, err = os.Open(args.samplesArchiveFilename)
	if err != nil {
		// err is guaranteed to be os.PathError
		pathErr := err.(*os.PathError)
		return errors.Errorf("failed to open samples archive '%s': %v",
			args.samplesArchiveFilename, pathErr.Err.Error())
	}

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Creator = client.Classifiers
	cmd.Deleter = client.Classifiers
	cmd.Trainer = client.Classifiers
	cmd.ClassifierName = args.classifierName
	return nil
}
