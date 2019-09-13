package main

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/ioutils"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"strings"
)

func main() {
	var (
		globalFlags = config.GlobalFlags{}

		app = kingpin.New("surf", "surf - the official command line client for waives.io.")

		listCmd   = app.Command("list", "List waives resources.")
		uploadCmd = app.Command("upload", "Upload waives resources.")
		deleteCmd = app.Command("delete", "Delete waives resources.")
		createCmd = app.Command("create", "Create waives resources.")

		ctx, canceller = context.WithCancel(context.Background())
	)

	go handleInterrupt(canceller)

	commands.ConfigureLoginCommand(ctx, app, &globalFlags)
	commands.ConfigureListModulesCommand(ctx, listCmd, &globalFlags)
	commands.ConfigureListClassifiersCmd(ctx, listCmd, &globalFlags)
	commands.ConfigureListExtractorsCmd(ctx, listCmd, &globalFlags)
	commands.ConfigureListDocumentsCmd(ctx, listCmd, &globalFlags)
	commands.ConfigureUploadExtractorCommand(ctx, uploadCmd, &globalFlags)
	commands.ConfigureDeleteExtractorCmd(ctx, deleteCmd, &globalFlags)
	commands.ConfigureDeleteClassifierCmd(ctx, deleteCmd, &globalFlags)
	commands.ConfigureDeleteDocumentCmd(ctx, deleteCmd, &globalFlags)
	commands.ConfigureCreateClassifierCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureCreateExtractorCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureCreateExtractorTemplateCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureCreateDocumentCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureReadCommand(ctx, app, &globalFlags)
	commands.ConfigureExtractCommand(ctx, app, &globalFlags)
	commands.ConfigureClassifyCommand(ctx, app, &globalFlags)
	commands.ConfigureUploadClassifierCommand(ctx, uploadCmd, &globalFlags)
	commands.ConfigureRedactWithExtractionCommand(ctx, app, &globalFlags)

	app.Flag("client-id", "Client ID").
		Short('i').
		PlaceHolder("id").
		StringVar(&globalFlags.ClientId)
	app.Flag("client-secret", "Client secret").
		Short('s').
		PlaceHolder("secret").
		StringVar(&globalFlags.ClientSecret)
	app.Flag("log-http", "Log HTTP requests and responses as they happen, "+
		"to a file.").
		PlaceHolder("file").
		OpenFileVar(&globalFlags.LogHttp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	app.Flag("version", "Show the application version.").
		PreAction(func(parseContext *kingpin.ParseContext) error {
			fmt.Println(ch360.Version)
			os.Exit(0)
			return nil
		}).
		Bool()

	app.UsageTemplate(kingpin.CompactUsageTemplate)
	app.HelpFlag.Hidden()

	defer ioutils.TryClose(globalFlags.LogHttp)

	_, err := app.Parse(os.Args[1:])
	exitOnErr(err)
}

func handleInterrupt(canceller context.CancelFunc) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan // ctrl-c received
	canceller()
}

func exitOnErr(err error) {
	if err != nil && errors.Cause(err) != context.Canceled {
		msg := "Error: " + err.Error()
		if !strings.HasSuffix(msg, ".") {
			msg = msg + "."
		}
		_, _ = fmt.Fprintln(os.Stderr, msg)

		os.Exit(1)
	}
}
