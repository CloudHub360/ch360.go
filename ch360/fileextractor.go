package ch360

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360/results"
	"io"
)

type FileExtractor struct {
	docCreator   DocumentCreator
	docExtractor DocumentExtractor
	docDeleter   DocumentDeleter
}

func NewFileExtractor(creator DocumentCreator, extractor DocumentExtractor, deleter DocumentDeleter) *FileExtractor {
	return &FileExtractor{
		docCreator:   creator,
		docExtractor: extractor,
		docDeleter:   deleter,
	}
}

func (f *FileExtractor) Extract(ctx context.Context, fileContents io.Reader, extractorName string) (*results.ExtractionResult, error) {
	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	document, err := f.docCreator.Create(context.Background(), fileContents)
	if err != nil {
		return nil, err
	}

	result, extractErr := f.docExtractor.Extract(ctx, document.Id, extractorName)

	if document.Id != "" {
		// Always delete the document, even if Extract returned an error.
		// Don't cancel on ctrl-c.
		err = f.docDeleter.Delete(context.Background(), document.Id)
	}

	// Return the extract err if we have one
	if extractErr != nil {
		return nil, extractErr
	}

	return result, nil
}
