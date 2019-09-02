package ch360

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360/request"
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

func (f *FileRedactor) Redact(ctx context.Context, fileContents io.Reader,
	extractorName string) (io.ReadCloser, error) {

	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	document, err := f.docCreator.Create(context.Background(), fileContents)
	if err != nil {
		return nil, err
	}

	defer func() {
		if document.Id != "" {
			// Always delete the document, even if extract / redact returned an error.
			// Don't cancel on ctrl-c.
			_ = f.docDeleter.Delete(context.Background(), document.Id)
		}
	}()

	// use the results from extraction...
	extractionResult, err := f.docExtractor.ExtractForRedaction(ctx, document.Id, extractorName)

	if err != nil {
		return nil, err
	}

	// ...as the input to redaction
	redactRequest := (*request.RedactedPdfRequest)(extractionResult)
	return f.docRedactor.Redact(ctx, document.Id, *redactRequest)
}
