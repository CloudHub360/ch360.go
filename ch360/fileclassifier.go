package ch360

import (
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

	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	document, err := f.docCreator.Create(context.Background(), fileContents)
	if err != nil {
		return nil, err
	}

	defer func() {
		if document.Id != "" {
			// Always delete the document, even if Classify returns an error.
			// Don't cancel on ctrl-c.
			f.docDeleter.Delete(context.Background(), document.Id)
		}
	}()

	return f.docClassifier.Classify(ctx, document.Id, classifierName)
}
