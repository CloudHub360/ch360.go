package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

//go:generate mockery -name "DocumentDeleterGetter"
type DocumentDeleterGetter interface {
	ch360.DocumentDeleter
	ch360.DocumentGetter
}

type deleteDocumentArgs struct {
	documentIds []string
	deleteAll   bool
}

type DeleteDocumentCmd struct {
	Client      DocumentDeleterGetter
	DocumentIDs []string
	DeleteAll   bool
}

// ConfigureDeleteDocumentCmd configures kingpin with the 'delete document' command.
func ConfigureDeleteDocumentCmd(ctx context.Context, deleteCmd *kingpin.CmdClause, flags *config.
	GlobalFlags) {
	args := &deleteDocumentArgs{}
	deleteDocumentCmd := &DeleteDocumentCmd{}

	deleteDocumentCli := deleteCmd.Command("document", "Delete waives documents.").
		Action(func(parseContext *kingpin.ParseContext) error {
			msg := fmt.Sprintf("Deleting %d documents... ", len(args.documentIds))
			if args.deleteAll {
				msg = "Deleting all documents... "
			}
			return ExecuteWithMessage(msg,
				func() error {
					err := deleteDocumentCmd.initFromArgs(args, flags)
					if err != nil {
						return err
					}
					return deleteDocumentCmd.Execute(ctx)
				})
		})

	deleteDocumentCli.
		Arg("ID", "The IDs of the document(s) to delete.").
		StringsVar(&args.documentIds)

	deleteDocumentCli.
		Flag("all", "Delete all documents.").
		BoolVar(&args.deleteAll)

	deleteDocumentCli.Validate(func(clause *kingpin.CmdClause) error {
		if !args.deleteAll && len(args.documentIds) == 0 {
			return errors.New("Please specify either --all or the document IDs to delete.")
		}

		if args.deleteAll && len(args.documentIds) > 0 {
			return errors.New("Please specify either --all or the document IDs to delete, " +
				"but not both.")
		}

		return nil
	})
}

// Execute is the entry point of the 'delete documents' command.
func (cmd *DeleteDocumentCmd) Execute(ctx context.Context) error {
	if cmd.DeleteAll {
		allDocIds, err := cmd.retrieveAllDocumentIds(ctx)
		if err != nil {
			return err
		}

		cmd.DocumentIDs = allDocIds
	}

	for _, docId := range cmd.DocumentIDs {
		err := cmd.Client.Delete(ctx, docId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *DeleteDocumentCmd) retrieveAllDocumentIds(ctx context.Context) ([]string, error) {
	allDocuments, err := cmd.Client.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var docIds []string
	for _, doc := range allDocuments {
		docIds = append(docIds, doc.Id)
	}
	return docIds, nil
}

func (cmd *DeleteDocumentCmd) initFromArgs(args *deleteDocumentArgs, flags *config.GlobalFlags) error {
	cmd.DocumentIDs = args.documentIds
	cmd.DeleteAll = args.deleteAll

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = client.Documents
	return nil
}
