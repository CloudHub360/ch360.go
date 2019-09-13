package commands

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

func addFileHandlingFlagsTo(globalFlags *config.GlobalFlags, cmdClause *kingpin.CmdClause) {

	cmdClause.Flag("multiple-files",
		"Write results output to multiple files with the same basename as the input").
		Short('m').
		BoolVar(&globalFlags.MultiFileOut)
	cmdClause.Flag("output-file", "Write all results to the specified file").
		Short('o').
		PlaceHolder("file").
		StringVar(&globalFlags.OutputFile)
	cmdClause.Flag("progress", "Show a progress bar (only for use with -o or -m).").
		Short('p').
		BoolVar(&globalFlags.ShowProgress)

	cmdClause.Validate(func(clause *kingpin.CmdClause) error {
		// Only show the progress bar if stdout is redirected, or -o or -m are used
		if globalFlags.ShowProgress && !globalFlags.CanShowProgressBar() {
			return errors.New("The --progress / -p option can only be used when " +
				"redirecting stdout, or in combination with -o or -m.")
		}
		return nil
	})
}
