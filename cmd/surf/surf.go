package main

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/docopt/docopt-go"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {
	usage := `surf - the official command line client for waives.io.

Usage:
  surf login [options]
  surf ` + new(commands.UploadClassifier).Usage() + ` <name> <classifier-file> [options]
  surf ` + new(commands.CreateClassifier).Usage() + ` <name> <samples-zip> [options]
  surf ` + new(commands.CreateExtractor).Usage() + ` <name> <config-file> [options]
  surf ` + new(commands.DeleteClassifier).Usage() + ` <name> [options]
  surf ` + new(commands.DeleteExtractor).Usage() + ` <name> [options]
  surf ` + new(commands.ListClassifiers).Usage() + ` [options]
  surf ` + new(commands.ListExtractors).Usage() + ` [options]
  surf ` + new(commands.ClassifyCommand).Usage() + ` <file> <classifier> [options]
  surf ` + new(commands.Extract).Usage() + ` <file> <extractor> [options]
  surf ` + new(commands.Read).Usage() + ` <file> (pdf|txt|wvdoc) [options]
  surf -h | --help
  surf -v | --version

Options:
  -h, --help                       : Show this help message
  -v, --version                    : Show application version
  -i, --client-id <id>             : Client ID
  -s, --client-secret <secret>     : Client secret
  -f, --output-format <format>     : Output format for classification and extraction results.        
                                     Allowed values: table, csv, json [default: table]
  -o, --output-file <file>         : Write all results to the specified file
  -m, --multiple-files             : Write results output to multiple files with the same
                                   : basename as the input
  -p, --progress                   : Show progress when classifying files. Only visible when
                                     redirecting stdout or in conjunction with -m or -o.
`

	filenameExamples := `
Filename and glob pattern examples:
  file1.pdf        : Specific file
  *.*              : All files in the current folder
  *.pdf            : All PDFs in the current folder
  foo/*.tif        : All TIFs in folder foo
  bar/**/*.*       : All files in subfolders of folder bar`

	// Replace slashes with OS-specific path separators
	usage = usage + filepath.FromSlash(filenameExamples)

	args, err := docopt.ParseArgs(usage, nil, ch360.Version)
	exitOnErr(err)

	runParams, err := config.NewRunParamsFromArgs(args)
	exitOnErr(err)

	ctx, canceller := context.WithCancel(context.Background())
	go handleInterrupt(canceller)

	appDir, err := config.NewAppDirectory()
	exitOnErr(err)

	var cmd commands.Command
	if login, _ := args.Bool("login"); login {
		tokenRetriever := ch360.NewTokenRetriever(commands.DefaultHttpClient, ch360.ApiAddress)
		cmd = commands.NewLoginFrom(runParams, os.Stdout, appDir, tokenRetriever)
	} else {
		cmd, err = commands.CommandFor(runParams)
	}
	exitOnErr(err)

	exitOnErr(cmd.Execute(ctx))
}

func handleInterrupt(canceller context.CancelFunc) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan // ctrl-c received
	canceller()
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
