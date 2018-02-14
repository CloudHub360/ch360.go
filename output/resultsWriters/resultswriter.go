package resultsWriters

//go:generate mockery -name ResultsWriter
type ResultsWriter interface {
	Start() error
	WriteResult(filename string, result interface{}) error
	Finish() error
}
