package commands

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/waives/surf/ch360"
	"github.com/waives/surf/config"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
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

	createDocumentCli := createCmd.Command("document",
		"Create one or more waives documents from files.").
		Alias("documents").
		Action(func(parseContext *kingpin.ParseContext) error {
			msg := "Creating document..."
			if len(args.documentPaths) > 1 {
				msg = "Creating documents... "
			}
			_, _ = fmt.Fprintln(os.Stderr, msg)
			err := createDocumentCmd.initFromArgs(args, flags)

			if err != nil {
				return err
			}

			return createDocumentCmd.Execute(ctx)
		})

	createDocumentCli.
		Arg("files", "The file(s) to create documents from.").
		Required().
		StringsVar(&args.documentPaths)
}

func (cmd *CreateDocumentCmd) Execute(ctx context.Context) error {
	table := NewTable(os.Stdout, []string{"File", "ID", "Size", "Type", "SHA256"})

	for _, documentPath := range cmd.DocumentPaths {
		doc, err := cmd.createFromFile(ctx, documentPath)
		if err != nil {
			return err
		}
		table.Append([]string{
			truncateStringLeft(documentPath, 40), doc.Id, strconv.Itoa(doc.Size),
			doc.FileType,
			doc.Sha256})
	}
	table.Render()

	return nil
}

func (cmd *CreateDocumentCmd) createFromFile(ctx context.Context,
	documentPath string) (*ch360.Document, error) {
	documentFile, err := os.Open(documentPath)
	if err != nil {
		pathErr := err.(*os.PathError)
		return nil, errors.Errorf("unable to create document from file '%s': %s", documentPath,
			pathErr.Err.Error())
	}
	defer documentFile.Close()

	doc, err := cmd.Creator.Create(ctx, documentFile)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create document from file '%s'",
			documentPath)
	}
	return &doc, nil
}

func (cmd *CreateDocumentCmd) initFromArgs(args *createDocumentArgs, flags *config.GlobalFlags) error {

	client, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Creator = client.Documents
	cmd.DocumentPaths, err = GlobMany(args.documentPaths)

	return err
}

func truncateStringLeft(str string, maxLength int) string {
	truncated := str
	strLen := len(str)
	if strLen > maxLength {
		if maxLength > 3 {
			maxLength -= 3
		}
		truncated = "..." + str[strLen-maxLength:strLen]
	}
	return truncated
}
