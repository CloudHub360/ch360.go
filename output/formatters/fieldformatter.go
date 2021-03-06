package formatters

import (
	"github.com/waives/surf/ch360/results"
	"strings"
)

// shouldOutputAlternatives is a helper method used to determine whether or not all
// alternative results should be output.
func (f FieldFormatter) shouldOutputAlternatives() bool {
	if !f.FieldResult.Rejected {
		return true
	}

	// Field is rejected
	return f.FieldResult.RejectReason == "MultipleResults"
}

// Results returns an array of results, taking field rejection and alternative results
// into account.
func (f FieldFormatter) Results() []string {
	results := make([]string, 0)
	if f.FieldResult.Result == nil {
		return results
	}

	results = append(results, f.FieldResult.Result.Text)

	if f.shouldOutputAlternatives() {
		for _, alternativeResult := range f.FieldResult.AlternativeResults {
			results = append(results, alternativeResult.Text)
		}
	}

	return results
}

// Returns a comma-separated joined string of Results(), or NoResultText if Results()
// returns an empty array.
func (f FieldFormatter) String() string {
	if len(f.Results()) == 0 {
		return f.NoResultStr
	}

	return strings.Join(f.Results(), f.Separator)
}

// FieldFormatter formats a results.FieldResult.
type FieldFormatter struct {
	FieldResult results.FieldResult
	Separator   string
	NoResultStr string
}

func NewFieldFormatter(result results.FieldResult, separator, noResultStr string) *FieldFormatter {
	return &FieldFormatter{
		FieldResult: result,
		Separator:   separator,
		NoResultStr: noResultStr,
	}
}
