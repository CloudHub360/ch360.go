package config

import (
	"github.com/mattn/go-isatty"
	"os"
)

type GlobalFlags struct {
	MultiFileOut bool
	OutputFile   string
	ShowProgress bool
	ClientId     string
	ClientSecret string
	LogHttp      *os.File
}

func (r *GlobalFlags) CanShowProgressBar() bool {
	return r.IsOutputSpecified() || IsOutputRedirected()
}

func IsOutputRedirected() bool {
	fd := os.Stdout.Fd()
	return !isatty.IsTerminal(fd) &&
		!isatty.IsCygwinTerminal(fd)
}

func (r *GlobalFlags) IsOutputSpecified() bool {
	return r.OutputFile != "" || r.MultiFileOut
}
