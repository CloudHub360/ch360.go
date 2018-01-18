package commands

import (
	"encoding/json"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/pkg/errors"
	"io"
	"path/filepath"
)

type JsonClassifyResultsWriter struct {
	writerProvider WriterProvider
	startCalled    bool
	writingStarted bool
}

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

func NewJsonClassifyResultsWriter(writerProvider WriterProvider) *JsonClassifyResultsWriter {
	return &JsonClassifyResultsWriter{
		writerProvider: writerProvider,
	}
}

func (writer *JsonClassifyResultsWriter) Start() {
	writer.startCalled = true
}

func (writer *JsonClassifyResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	if !writer.startCalled {
		return errors.New("Start() must be called before WriteResult()")
	}

	out, err := writer.writerProvider.Provide(filename)

	if err != nil {
		return err
	}

	if writer.writingStarted {
		fmt.Fprint(out, ",")
		fmt.Fprintln(out, "")
	} else {
		fmt.Fprint(out, "[")
		writer.writingStarted = true
	}

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

	_, err = out.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (writer *JsonClassifyResultsWriter) Finish() {
	out, _ := writer.writerProvider.Provide("")

	if out == nil {
		return
	}

	if writer.writingStarted {
		fmt.Fprint(out, "]")
	}
}
