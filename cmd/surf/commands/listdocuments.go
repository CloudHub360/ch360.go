package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strconv"
)

type ListDocumentsCmd struct {
	Client ch360.DocumentGetter
}

// Configures kingpin with the 'list documents' command
func ConfigureListDocumentsCmd(ctx context.Context, parentCmd *kingpin.CmdClause, flags *config.GlobalFlags) {
	listDocumentsCmd := &ListDocumentsCmd{}
	parentCmd.Command("documents", "List all available documents.").
		Alias("document").
		Action(func(parseContext *kingpin.ParseContext) error {
			err := listDocumentsCmd.initFromArgs(flags)
			if err != nil {
				return err
			}
			return listDocumentsCmd.Execute(ctx)
		})
}

// Executes the command.
func (cmd *ListDocumentsCmd) Execute(ctx context.Context) error {
	documents, err := cmd.Client.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(documents) == 0 {
		fmt.Println("No documents found.")
		return nil
	}

	table := NewTable(os.Stdout, []string{"ID", "Size", "Type", "SHA256"})

	for _, document := range documents {
		table.Append([]string{document.Id, strconv.Itoa(document.Size), document.FileType,
			document.Sha256})
	}
	table.Render()

	return nil
}

func (cmd *ListDocumentsCmd) initFromArgs(flags *config.GlobalFlags) error {
	var err error
	apiClient, err := initApiClient(flags.ClientId, flags.ClientSecret, flags.LogHttp)

	if err != nil {
		return err
	}

	cmd.Client = apiClient.Documents
	return nil
}
