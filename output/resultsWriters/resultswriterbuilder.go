package resultsWriters

import (
	"os"

	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
	"github.com/spf13/afero"
)

func NewResultsWriterFor(params *config.GlobalFlags, fileExtension string,
	verb config.Verb) (ResultsWriter, error) {

	formatter, err := formatters.NewResultsFormatterFor(params, verb)

	if err != nil {
		return nil, err
	}

	if params.MultiFileOut {
		// Write output to a file "next to" each input file, with the specified extension
		writerFactory := sinks.NewExtensionSwappingFileSinkFactory(fileExtension)
		return NewIndividualResultsWriter(writerFactory, formatter), nil
	} else if params.OutputFile != "" {
		// Write output to a single "combined results" file with the specified filename
		return NewCombinedResultsWriter(sinks.NewBasicFileSink(afero.NewOsFs(), params.OutputFile), formatter), nil
	} else {
		// Write output to the console
		return NewCombinedResultsWriter(sinks.NewBasicWriterSink(os.Stdout), formatter), nil
	}
}
