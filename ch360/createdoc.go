package ch360

import (
	"context"
	"io"
)

// CreateDocumentFor is a helper function that creates a document with a specified
// DocumentCreator, runs a provided function with the newly-created document,
// and then deletes the document. It does not accept a context.Context as it assumes
// the user would want neither the creation nor deletion to be interrupted,
// since that could lead to 'leaked' documents in waives.
func CreateDocumentFor(contents io.Reader, creator DocumentCreator,
	deleter DocumentDeleter,
	fn func(Document) error) error {

	ctx := context.Background()

	document, err := creator.Create(ctx, contents)
	if err != nil {
		return err
	}

	defer func() {
		// always delete the document, even if fn fails.
		if document.Id != "" {
			deleter.Delete(ctx, document.Id)
		}
	}()

	return fn(document)
}
