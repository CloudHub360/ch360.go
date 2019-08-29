package ch360

import (
	"context"
	"errors"
)

var TotalDocumentSlots = 30

var ErrDocSlotsFull = errors.New("all document slots are full")

// GetFreeDocSlots is a helper function to retrieve the number of available document slots in
// waives, and return an error if there are none.
func GetFreeDocSlots(ctx context.Context, getter DocumentGetter, totalSlots int) (int,
	error) {
	documentList, err := getter.GetAll(ctx)

	if err != nil {
		return 0, err
	}

	slots := totalSlots - len(documentList)

	// since the total doc slots is currently hardcoded,
	// it's possible (even likely) that the actual number
	// of documents in waives can be be greater than TotalDocumentSlots.
	if slots <= 0 {
		return slots, ErrDocSlotsFull
	}

	return slots, nil
}
