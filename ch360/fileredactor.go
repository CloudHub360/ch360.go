package ch360

import (
	"context"
	"github.com/waives/surf/ch360/request"
	"io"
)

type FileRedactor struct {
	docCreator   DocumentCreator
	docExtractor DocumentExtractor
	docRedactor  DocumentRedactor
	docDeleter   DocumentDeleter
}

func NewFileRedactor(creator DocumentCreator, extractor DocumentExtractor,
	redactor DocumentRedactor,
	deleter DocumentDeleter) *FileRedactor {
	return &FileRedactor{
		docCreator:   creator,
		docExtractor: extractor,
		docRedactor:  redactor,
		docDeleter:   deleter,
	}
}

// Redact creates a document from the fileContents, performs extraction on it,
// then redacts it with the results of the extraction.
func (f *FileRedactor) Redact(ctx context.Context, fileContents io.Reader,
	extractorName string) (io.ReadCloser, error) {
	var (
		redacted io.ReadCloser
		err      error
	)

	err = CreateDocumentFor(fileContents, f.docCreator, f.docDeleter,
		func(document Document) error {
			// use the results from extraction...
			extractionResult, err := f.docExtractor.ExtractForRedaction(ctx, document.Id,
				extractorName)

			if err != nil {
				return err
			}

			// ...as the input to redaction
			redactRequest := (*request.RedactedPdfRequest)(extractionResult)
			redacted, err = f.docRedactor.Redact(ctx, document.Id, *redactRequest)

			return err
		})

	return redacted, err
}
