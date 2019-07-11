package ioutils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestDrainClose(t *testing.T) {
	fixtures := []struct {
		reader        io.Reader
		expectedBytes []byte
	}{
		{
			reader:        nil,
			expectedBytes: []byte(nil),
		},
		{
			reader:        bytes.NewBufferString("some data"),
			expectedBytes: []byte("some data"),
		}, {
			reader:        bytes.NewReader([]byte("some data")),
			expectedBytes: []byte("some data"),
		},
	}
	for _, fixture := range fixtures {
		buffer, err := DrainClose(fixture.reader)
		actualBytes := buffer.Bytes()

		assert.NoError(t, err)
		assert.Equal(t, fixture.expectedBytes, actualBytes)
	}
}
