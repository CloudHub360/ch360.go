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

//go:generate stringer -type=Verb
type Verb int

const (
	Extract Verb = iota
	Classify
	Login
	List
	Create
	Delete
	Upload
	Read
)

//go:generate stringer -type=Noun
type Noun int

const (
	Classifier Noun = iota
	Extractor
	Module
	ExtractorTemplate
)

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
