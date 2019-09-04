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

// Read creates a document from fileContents, performs a read,
// then returns the read results in the format according to mode.
func (f *FileReader) Read(ctx context.Context, fileContents io.Reader, mode ReadMode) (io.ReadCloser, error) {
	var (
		result io.ReadCloser
		err    error
	)

	err = CreateDocumentFor(fileContents, f.docCreator, f.docDeleter,
		func(document Document) error {
			if err = f.docReader.Read(ctx, document.Id); err == nil {
				result, err = f.docReader.ReadResult(ctx, document.Id, mode)
			}

			return err
		})

	return result, err
}
