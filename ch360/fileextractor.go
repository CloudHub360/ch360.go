package ch360

import (
	"context"
	"github.com/waives/surf/ch360/results"
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

// Extract creates a document, performs extraction, deletes the doc,
// then returns the extraction result.
func (f *FileExtractor) Extract(ctx context.Context, fileContents io.Reader, extractorName string) (*results.ExtractionResult, error) {
	var (
		extractionResult *results.ExtractionResult
		err              error
	)

	err = CreateDocumentFor(fileContents, f.docCreator, f.docDeleter,
		func(document Document) error {
			extractionResult, err = f.docExtractor.Extract(ctx, document.Id, extractorName)
			return err

		})

	return extractionResult, err
}
