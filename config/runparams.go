package config

import (
	"errors"
	"github.com/docopt/docopt-go"
	"github.com/mattn/go-isatty"
	"os"
)

type RunParams struct {
	Login    bool
	Classify bool
	Extract  bool
	Create   bool
	Delete   bool
	List     bool
	Read     bool

	Extractor   bool
	Classifier  bool
	Extractors  bool
	Classifiers bool

	ClassifierName string `docopt:"<classifier>"`
	ExtractorName  string `docopt:"<extractor>"`

	MultiFileOut bool   `docopt:"-m,--multiple-files"`
	OutputFile   string `docopt:"-o,--output-file"`
	OutputFormat string `docopt:"-f,--output-format"`
	ShowProgress bool   `docopt:"-p,--progress"`
	ClientId     string `docopt:"-i,--client-id"`
	ClientSecret string `docopt:"-s,--client-secret"`

	ConfigPath  string `docopt:"<config-file>"`
	SamplesPath string `docopt:"<samples-zip>"`
	FilePattern string `docopt:"<file>"`
	Name        string `docopt:"<name>"`

	ReadPDF   bool `docopt:"pdf"`
	ReadText  bool `docopt:"txt"`
	ReadWvdoc bool `docopt:"wvdoc"`

	args docopt.Opts
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
	Read
)

//go:generate stringer -type=Noun
type Noun int

const (
	Classifier Noun = iota
	Extractor
)

func (r RunParams) Verb() Verb {
	if r.Create {
		return Create
	} else if r.Delete {
		return Delete
	} else if r.List {
		return List
	} else if r.Login {
		return Login
	} else if r.Classify {
		return Classify
	} else if r.Extract {
		return Extract
	} else if r.Read {
		return Read
	}

	return -1
}

func NewRunParamsFromArgs(args docopt.Opts) (*RunParams, error) {
	runParams := RunParams{}

	err := args.Bind(&runParams)

	if err != nil {
		return nil, err
	}

	err = runParams.Validate()
	if err != nil {
		return nil, err
	}

	runParams.args = args

	// Only show progress if stdout is being redirected
	if !shouldShowProgressBar(runParams.MultiFileOut || runParams.OutputFile != "") {
		runParams.ShowProgress = false
	}

	return &runParams, nil
}

func (r RunParams) Validate() error {
	if r.MultiFileOut && r.OutputFile != "" {
		return errors.New("The --multiple-files (-m) and --output-file (-o) options cannot be used together.")
	}

	return nil
}

// Returns the noun (Classifier or Extractor) specified, or -1 if not present.
func (r RunParams) Noun() Noun {
	if r.Extractor || r.Extractors {
		return Extractor
	} else if r.Classifier || r.Classifiers {
		return Classifier
	}

	return -1
}

func (r RunParams) Args() docopt.Opts {
	return r.args
}

func shouldShowProgressBar(writingToFile bool) bool {
	return writingToFile || isRedirected(os.Stdout.Fd())
}

func isRedirected(fd uintptr) bool {
	return !isatty.IsTerminal(fd) &&
		!isatty.IsCygwinTerminal(fd)
}
