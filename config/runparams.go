package config

import (
	"github.com/docopt/docopt-go"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"os"
)

type RunParams struct {
	Login    bool
	Classify bool
	Extract  bool
	Create   bool
	Upload   bool
	Delete   bool
	List     bool
	Read     bool

	Extractor         bool
	Classifier        bool
	Extractors        bool
	ExtractorTemplate bool `docopt:"extractor-template"`
	Classifiers       bool
	Modules           bool
	Module            bool

	ClassifierName string `docopt:"<classifier>"`
	ExtractorName  string `docopt:"<extractor>"`

	MultiFileOut bool   `docopt:"-m,--multiple-files"`
	OutputFile   string `docopt:"-o,--output-file"`
	OutputFormat string `docopt:"-f,--output-format"`
	ShowProgress bool   `docopt:"-p,--progress"`
	ClientId     string `docopt:"-i,--client-id"`
	ClientSecret string `docopt:"-s,--client-secret"`
	LogHttp      string `docopt:"--log-http"`

	ModulesTemplate string   `docopt:"-t,--from-template"`
	ModuleIds       []string `docopt:"<module-ids>"`

	ConfigPath     string `docopt:"<config-file>"`
	SamplesPath    string `docopt:"<samples-zip>"`
	ClassifierPath string `docopt:"<classifier-file>"`
	FilePattern    string `docopt:"<file>"`
	Name           string `docopt:"<name>"`

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
	} else if r.Upload {
		return Upload
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
	} else if r.ExtractorTemplate {
		return ExtractorTemplate
	} else if r.Classifier || r.Classifiers {
		return Classifier
	} else if r.Module || r.Modules {
		return Module
	}

	return -1
}

func (r RunParams) Args() docopt.Opts {
	return r.args
}

func shouldShowProgressBar(writingToFile bool) bool {
	return writingToFile || IsOutputRedirected()
}

func IsOutputRedirected() bool {
	fd := os.Stdout.Fd()
	return !isatty.IsTerminal(fd) &&
		!isatty.IsCygwinTerminal(fd)
}

func (r *RunParams) IsOutputSpecified() bool {
	return r.OutputFile != "" || r.MultiFileOut
}
