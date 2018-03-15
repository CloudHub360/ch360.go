package resultsWriters

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/CloudHub360/ch360.go/output/sinks"
	"github.com/spf13/afero"
	"os"
)

func NewResultsWriterFor(params *config.RunParams) (ResultsWriter, error) {

	formatter, err := formatters.NewResultsFormatterFor(params)

	if err != nil {
		return nil, err
	}

	if params.MultiFileOut {
		// Write output to a file "next to" each input file, with the specified extension
		outputFileExtension := fileExtensionFor(params)
		writerFactory := sinks.NewExtensionSwappingFileSinkFactory(outputFileExtension)
		return NewIndividualResultsWriter(writerFactory, formatter), nil
	} else if params.OutputFile != "" {
		// Write output to a single "combined results" file with the specified filename
		return NewCombinedResultsWriter(sinks.NewBasicFileSink(afero.NewOsFs(), params.OutputFile), formatter), nil
	} else {
		// Write output to the console
		return NewCombinedResultsWriter(sinks.NewBasicWriterSink(os.Stdout), formatter), nil
	}
}

func fileExtensionFor(params *config.RunParams) string {
	if params.Verb() == config.Read {
		if params.ReadPDF {
			return ".ocr.pdf"
		} else {
			return ".ocr.txt"
		}
	}

	var outputFormatExtensions = map[formatters.OutputFormat]string{
		formatters.Table: ".tab",
		formatters.Json:  ".json",
		formatters.Csv:   ".csv",
	}
	return outputFormatExtensions[formatters.OutputFormat(params.OutputFormat)]
}
