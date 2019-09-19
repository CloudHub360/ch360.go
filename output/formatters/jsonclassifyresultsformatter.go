package formatters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/waives/surf/ch360/results"
	"io"
	"path/filepath"
)

type JsonClassifyResultsFormatter struct {
	resultsWritten bool
	headerWritten  bool
}

var _ ResultsFormatter = (*JsonClassifyResultsFormatter)(nil)

type classifyDocumentOutput struct {
	Filename string                       `json:"filename"`
	Results  classifyDocumentResultOutput `json:"classification_results"`
}

type classifyDocumentResultOutput struct {
	DocumentType       string                                    `json:"document_type"`
	IsConfident        bool                                      `json:"is_confident"`
	RelativeConfidence float64                                   `json:"relative_confidence"`
	Scores             []classifyDocumentResultDocumentTypeScore `json:"document_type_scores"`
}

type classifyDocumentResultDocumentTypeScore struct {
	DocumentType string  `json:"document_type"`
	Score        float64 `json:"score"`
}

func NewJsonClassifyResultsFormatter() *JsonClassifyResultsFormatter {
	return &JsonClassifyResultsFormatter{}
}

func (f *JsonClassifyResultsFormatter) WriteResult(writer io.Writer, filename string, result interface{}, options FormatOption) error {

	classificationResult, ok := result.(*results.ClassificationResult)

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected type: %T", result))
	}

	if options&IncludeHeader == IncludeHeader {
		f.headerWritten = true
		fmt.Fprint(writer, "[") // header
	} else if f.resultsWritten {
		fmt.Fprint(writer, ",\n") // write separator
	}

	var scores []classifyDocumentResultDocumentTypeScore
	for _, score := range classificationResult.DocumentTypeScores {
		scores = append(scores, classifyDocumentResultDocumentTypeScore{DocumentType: score.DocumentType, Score: score.Score})
	}

	output := &classifyDocumentOutput{
		Filename: filepath.FromSlash(filename),
		Results: classifyDocumentResultOutput{
			DocumentType:       classificationResult.DocumentType,
			IsConfident:        classificationResult.IsConfident,
			RelativeConfidence: classificationResult.RelativeConfidence,
			Scores:             scores,
		},
	}

	bytes, err := json.MarshalIndent(output, "", "  ")

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)

	f.resultsWritten = true

	return err
}

func (f *JsonClassifyResultsFormatter) Flush(writer io.Writer) error {
	if f.headerWritten {
		_, err := fmt.Fprint(writer, "]")
		return err
	}
	return nil
}
