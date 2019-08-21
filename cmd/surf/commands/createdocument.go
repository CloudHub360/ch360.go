package commands

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type CreateDocumentCmd struct {
	Creator       ch360.DocumentCreator
	DocumentPaths []string
}

type createDocumentArgs struct {
	documentPaths []string
}

func ConfigureCreateDocumentCmd(ctx context.Context, createCmd *kingpin.CmdClause,
	flags *config.GlobalFlags) {
	args := &createDocumentArgs{}
	createDocumentCmd := &CreateDocumentCmd{}

	createDocumentCli := createCmd.Command("document", "Create waives document from a file.").
		Alias("documents").
		Action(func(parseContext *kingpin.ParseContext) error {
			msg := "Creating document... "
			if len(args.documentPaths) > 1 {
				msg = "Creating documents... "
			}
			return ExecuteWithMessage(msg,
				func() error {
					err := createDocumentCmd.initFromArgs(args, flags)

					if err != nil {
						return err
					}

					return createDocumentCmd.Execute(ctx)
				})
		})

	createDocumentCli.
		Arg("documents", "The file(s) to create documents from.").
		Required().
		StringsVar(&args.documentPaths)
}

func (cmd *CreateDocumentCmd) Execute(ctx context.Context) error {
	for _, documentPath := range cmd.DocumentPaths {
		err := cmd.createFromFile(ctx, documentPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *CreateDocumentCmd) createFromFile(ctx context.Context, documentPath string) error {
	documentFile, err := os.Open(documentPath)
	if err != nil {
		pathErr := err.(*os.PathError)
		return errors.Errorf("Unable to create document from file '%s': %s", documentPath, pathErr.Err.Error())
	}
	defer documentFile.Close()

	_, err = cmd.Creator.Create(ctx, documentFile)
	if err != nil {
		return errors.Wrapf(err, "Unable to create document from file '%s'",
			documentPath)
	}
	return nil
}

func (cmd *CreateDocumentCmd) initFromArgs(args *createDocumentArgs, flags *config.GlobalFlags) error {

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Creator = client.Documents
	cmd.DocumentPaths = args.documentPaths

	return nil
}
