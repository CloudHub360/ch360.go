package commands

import (
	"encoding/json"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/pkg/errors"
	"io"
)

type JsonClassifyResultsWriter struct {
	underlyingWriter io.Writer
	startCalled      bool
	writingStarted   bool
}

type classifyDocumentOutput struct {
	Filename string                       `json:"filename"`
	Results  classifyDocumentResultOutput `json:"classification_results"`
}

type classifyDocumentResultOutput struct {
	DocumentType       string  `json:"document_type"`
	IsConfident        bool    `json:"is_confident"`
	RelativeConfidence float64 `json:"relative_confidence"`
}

func NewJsonClassifyResultsWriter(writer io.Writer) *JsonClassifyResultsWriter {
	return &JsonClassifyResultsWriter{
		underlyingWriter: writer,
	}
}

func (writer *JsonClassifyResultsWriter) Start() {
	writer.startCalled = true
}

func (writer *JsonClassifyResultsWriter) WriteResult(filename string, result *types.ClassificationResult) error {
	if !writer.startCalled {
		return errors.New("Start() must be called before WriteResult()")
	}

	if writer.writingStarted {
		fmt.Fprintln(writer.underlyingWriter, "")
	} else {
		fmt.Fprint(writer.underlyingWriter, "[")
		writer.writingStarted = true
	}

	output := &classifyDocumentOutput{
		Filename: filename,
		Results: classifyDocumentResultOutput{
			DocumentType:       result.DocumentType,
			IsConfident:        result.IsConfident,
			RelativeConfidence: result.RelativeConfidence,
		},
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return err
	}

	writer.underlyingWriter.Write(bytes)
	return nil
}

func (writer *JsonClassifyResultsWriter) Finish() {
	fmt.Fprintln(writer.underlyingWriter, "]")
}
