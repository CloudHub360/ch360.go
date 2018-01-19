package formatters

import (
	"encoding/json"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"io"
	"path/filepath"
)

type JsonClassifyResultsFormatter struct{}

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

func (f *JsonClassifyResultsFormatter) WriteHeader(writer io.Writer) error {
	fmt.Fprint(writer, "[")
	return nil
}

func (f *JsonClassifyResultsFormatter) WriteResult(writer io.Writer, filename string, result *types.ClassificationResult) error {
	var scores []classifyDocumentResultDocumentTypeScore
	for _, score := range result.DocumentTypeScores {
		scores = append(scores, classifyDocumentResultDocumentTypeScore{DocumentType: score.DocumentType, Score: score.Score})
	}

	output := &classifyDocumentOutput{
		Filename: filepath.FromSlash(filename),
		Results: classifyDocumentResultOutput{
			DocumentType:       result.DocumentType,
			IsConfident:        result.IsConfident,
			RelativeConfidence: result.RelativeConfidence,
			Scores:             scores,
		},
	}

	bytes, err := json.MarshalIndent(output, "", "  ")

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)

	return err
}

func (f *JsonClassifyResultsFormatter) WriteSeparator(writer io.Writer) error {
	_, err := fmt.Fprint(writer, ",\n")
	return err
}

func (f *JsonClassifyResultsFormatter) WriteFooter(writer io.Writer) error {
	_, err := fmt.Fprint(writer, "]")
	return err
}
