package tests

import (
	"github.com/CloudHub360/ch360.go/ch360/results"
	"github.com/CloudHub360/ch360.go/output/formatters"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FieldFormatter(t *testing.T) {
	fixtures := []struct {
		fieldResult     results.FieldResult
		expectedResults []string
		expectedString  string
		separator       string
		noResultStr     string
	}{
		{
			fieldResult:     aFieldResult(false, []string{"result1", "alternative1", "alternative2"}, ""),
			expectedResults: []string{"result1", "alternative1", "alternative2"},
			expectedString:  "result1, alternative1, alternative2",
			separator:       ", ",
		}, {
			fieldResult:     aFieldResult(true, []string{"result1", "alternative1", "alternative2"}, "MultipleResults"),
			expectedResults: []string{"result1", "alternative1", "alternative2"},
			expectedString:  "result1, alternative1, alternative2",
			separator:       ", ",
		}, {
			fieldResult:     aFieldResult(true, []string{"result1", "alternative1", "alternative2"}, ""),
			expectedResults: []string{"result1"},
			expectedString:  "result1",
			separator:       ", ",
		}, {
			fieldResult:     aFieldResult(false, []string{"result1", "alternative1", "alternative2"}, ""),
			expectedResults: []string{"result1", "alternative1", "alternative2"},
			expectedString:  "result1|alternative1|alternative2",
			separator:       "|",
		}, {
			fieldResult:     aFieldResult(false, []string{}, ""),
			expectedResults: []string{},
			expectedString:  "(no result)",
			noResultStr:     "(no result)",
		}, {
			fieldResult:     aFieldResult(false, []string{}, ""),
			expectedResults: []string{},
			expectedString:  "",
			noResultStr:     "",
		},
	}

	for _, fixture := range fixtures {
		sut := formatters.NewFieldFormatter(fixture.fieldResult, fixture.separator, fixture.noResultStr)

		actualString := sut.String()
		actualResults := sut.Results()

		assert.Equal(t, fixture.expectedString, actualString)
		assert.Equal(t, fixture.expectedResults, actualResults)
	}
}

func aFieldResult(rejected bool, resultStrings []string, rejectReason string) results.FieldResult {
	fieldResult := results.FieldResult{
		Rejected:     rejected,
		RejectReason: rejectReason,
	}

	if len(resultStrings) == 0 {
		return fieldResult
	}

	fieldResult.Result = &results.InnerResult{
		Rejected: rejected,
		Text:     resultStrings[0],
	}

	for _, alternativeResult := range resultStrings[1:] {
		fieldResult.AlternativeResults = append(fieldResult.AlternativeResults, &results.InnerResult{
			Rejected: rejected,
			Text:     alternativeResult,
		})
	}

	return fieldResult
}
