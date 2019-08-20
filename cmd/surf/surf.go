package main

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/ioutils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
)

func main() {
	//	usage := `surf - the official command line client for waives.io.
	//
	//Usage:
	//  surf login [options]
	//  surf ` + new(commands.ClassifyCommand).Usage() + ` <file> <classifier> [options]
	//  surf ` + new(commands.Extract).Usage() + ` <file> <extractor> [options]
	//  surf ` + new(commands.Read).Usage() + ` <file> (pdf|txt|wvdoc) [options]
	//  surf -h | --help
	//  surf -v | --version
	//
	//Options:
	//  -h, --help                       : Show this help message
	//  -v, --version                    : Show application version
	//  -i, --client-id <id>             : Client ID
	//  -s, --client-secret <secret>     : Client secret
	//  -f, --output-format <format>     : Output format for classification and extraction results.
	//                                     Allowed values: table, csv, json [default: table]
	//  -o, --output-file <file>         : Write all results to the specified file
	//  -m, --multiple-files             : Write results output to multiple files with the same
	//                                   : basename as the input
	//  -p, --progress                   : Show progress when classifying files. Only visible when
	//                                     redirecting stdout or in conjunction with -m or -o.
	//  -t, --from-template <template>   : The extractor modules template to use when creating an
	//                                     extractor from modules.
	//  --log-http <file>                : Log HTTP requests and responses as they happen, to a file.
	//`
	//
	//	filenameExamples := `
	//Filename and glob pattern examples:
	//  file1.pdf        : Specific file
	//  *.*              : All files in the current folder
	//  *.pdf            : All PDFs in the current folder
	//  foo/*.tif        : All TIFs in folder foo
	//  bar/**/*.*       : All files in subfolders of folder bar`
	//
	//	// Replace slashes with OS-specific path separators
	//	usage = usage + filepath.FromSlash(filenameExamples)

	var (
		globalFlags = config.GlobalFlags{}

		app = kingpin.New("surf", "surf - the official command line client for waives.io.").
			Version(ch360.Version)

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
	commands.ConfigureUploadExtractorCommand(ctx, uploadCmd, &globalFlags)
	commands.ConfigureDeleteExtractorCmd(ctx, deleteCmd, &globalFlags)
	commands.ConfigureDeleteClassifierCmd(ctx, deleteCmd, &globalFlags)
	commands.ConfigureCreateClassifierCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureCreateExtractorCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureCreateExtractorTemplateCmd(ctx, createCmd, &globalFlags)
	commands.ConfigureReadCommand(ctx, app, &globalFlags)
	commands.ConfigureExtractCommand(ctx, app, &globalFlags)
	commands.ConfigureClassifyCommand(ctx, app, &globalFlags)
	commands.ConfigureUploadClassifierCommand(ctx, uploadCmd, &globalFlags)

	app.Flag("client-id", "Client ID").Short('i').
		Short('i').
		StringVar(&globalFlags.ClientId)
	app.Flag("client-secret", "Client secret").Short('s').
		Short('s').
		StringVar(&globalFlags.ClientSecret)
	app.Flag("log-http", "Log HTTP requests and responses as they happen, "+
		"to a file.").
		OpenFileVar(&globalFlags.LogHttp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	app.Flag("multiple-files",
		"Write results output to multiple files with the same basename as the input").
		Short('m').
		BoolVar(&globalFlags.MultiFileOut)
	app.Flag("output-file", "Write all results to the specified file").
		Short('o').
		StringVar(&globalFlags.OutputFile)

	defer ioutils.TryClose(globalFlags.LogHttp)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func handleInterrupt(canceller context.CancelFunc) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan // ctrl-c received
	canceller()
}
