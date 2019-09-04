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

	var (
		result *results.ClassificationResult
		err    error
	)
	err = CreateDocumentFor(fileContents, f.docCreator, f.docDeleter,
		func(document Document) error {
			result, err = f.docClassifier.Classify(ctx, document.Id, classifierName)
			return err
		})

	return result, err
}
