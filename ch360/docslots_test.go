package ch360_test

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_GetFreeDocSlots(t *testing.T) {
	fixtures := []struct {
		totalSlots    int
		presentDocs   ch360.DocumentList
		expectedSlots int
		expectedErr   error
		ctx           context.Context
	}{
		{
			totalSlots:    10,
			presentDocs:   nil,
			expectedSlots: 10,
			expectedErr:   nil,
			ctx:           context.Background(),
		}, {
			totalSlots:    3,
			presentDocs:   aListOfDocuments("1", "2", "3"),
			expectedSlots: 0,
			expectedErr:   ch360.ErrDocSlotsFull,
			ctx:           context.Background(),
		}, {
			totalSlots:  2,
			presentDocs: aListOfDocuments("1", "2", "3"),
			// it's expected that the slot count could be <0,
			// since the total slots is currently hardcoded to ch360.TotalDocumentSlots (30).
			expectedSlots: -1,
			expectedErr:   ch360.ErrDocSlotsFull,
			ctx:           context.Background(),
		}, {
			totalSlots:    10,
			presentDocs:   aListOfDocuments("1", "2", "3"),
			expectedSlots: 7,
			expectedErr:   nil,
			ctx:           context.Background(),
		},
	}

	docGetter := &mocks.DocumentGetter{}

	for _, fixture := range fixtures {
		docGetter.ExpectedCalls = nil
		docGetter.On("GetAll", mock.Anything).Return(fixture.presentDocs, nil)

		actualSlots, actualErr := ch360.GetFreeDocSlots(fixture.ctx, docGetter, fixture.totalSlots)

		assert.Equal(t, fixture.expectedSlots, actualSlots)
		assert.Equal(t, fixture.expectedErr, actualErr)
	}
}

func Test_GetFreeDocSlots_Returns_Err_From_DocGetter(t *testing.T) {
	expectedErr := errors.New("simulated error")
	docGetter := &mocks.DocumentGetter{}
	docGetter.ExpectedCalls = nil
	docGetter.On("GetAll", mock.Anything).Return(nil, expectedErr)

	_, actualErr := ch360.GetFreeDocSlots(context.Background(), docGetter, 10)

	assert.Equal(t, expectedErr, actualErr)
}

func aListOfDocuments(ids ...string) ch360.DocumentList {
	expected := make(ch360.DocumentList, len(ids))

	for index, id := range ids {
		expected[index] = ch360.Document{
			Id:       id,
			Size:     generators.Int(),
			Sha256:   generators.String("sha"),
			FileType: generators.String("fileType"),
		}
	}

	return expected
}
