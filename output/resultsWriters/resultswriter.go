package resultsWriters

import "github.com/CloudHub360/ch360.go/ch360/types"

//go:generate mockery -name ResultsWriter
type ResultsWriter interface {
	Start() error
	WriteResult(filename string, result *types.ClassificationResult) error
	Finish() error
}
