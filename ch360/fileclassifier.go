package ch360

import (
	"bytes"
	"context"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"io"
)

type FileClassifier struct {
	docCreator    DocumentCreator
	docClassifier DocumentClassifier
	docDeleter    DocumentDeleter
}

func NewFileClassifier(creator DocumentCreator, classifier DocumentClassifier,
	deleter DocumentDeleter) *FileClassifier {
	return &FileClassifier{
		docCreator:    creator,
		docClassifier: classifier,
		docDeleter:    deleter,
	}
}

func (f *FileClassifier) Classify(ctx context.Context, fileContents io.Reader,
	classifierName string) (*results.ClassificationResult, error) {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(fileContents)
	if err != nil {
		return nil, err
	}

	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	documentId, err := f.docCreator.Create(context.Background(), buf.Bytes())
	if err != nil {
		return nil, err
	}

	defer func() {
		if documentId != "" {
			// Always delete the document, even if Classify returns an error.
			// Don't cancel on ctrl-c.
			f.docDeleter.Delete(context.Background(), documentId)
		}
	}()

	return f.docClassifier.Classify(ctx, documentId, classifierName)
}
