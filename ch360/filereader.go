package ch360

import (
	"context"
	"io"
)

// Helper struct which creates a document from a file, performs a read, downloads the read
// result, then deletes the document.
type FileReader struct {
	docCreator DocumentCreator
	docReader  DocumentReader
	docDeleter DocumentDeleter
}

func NewFileReader(creator DocumentCreator, reader DocumentReader, deleter DocumentDeleter) *FileReader {
	return &FileReader{
		docCreator: creator,
		docReader:  reader,
		docDeleter: deleter,
	}
}

func (f *FileReader) Read(ctx context.Context, fileContents io.Reader, mode ReadMode) (io.ReadCloser, error) {
	// Use a different context here so we don't cancel this req on ctrl-c. We need
	// the docId result to perform cleanup
	documentId, err := f.docCreator.Create(context.Background(), fileContents)
	if err != nil {
		return nil, err
	}

	var (
		result  io.ReadCloser
		readErr error
	)
	if readErr = f.docReader.Read(ctx, documentId); readErr == nil {
		result, readErr = f.docReader.ReadResult(ctx, documentId, mode)
	}

	// Always delete the document, even if Read returned an error.
	// Don't cancel on ctrl-c.
	err = f.docDeleter.Delete(context.Background(), documentId)

	// Return the read err if we have one
	if readErr != nil {
		return nil, readErr
	}

	return result, err
}
