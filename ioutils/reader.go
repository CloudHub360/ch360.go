package ioutils

import (
	"bytes"
	"io"
)

// DrainClose will read all bytes of an io.Reader into a bytes.Buffer,
// to allow it to be used multiple times.
//
// If the reader is a *bytes.Buffer, a new bytes.Buffer will be returned
// from its backing []byte slice.
//
// It will also attempt to Close any provided
// reader, if it implements the io.Closer iface.
func DrainClose(reader io.Reader) (*bytes.Buffer, error) {
	if reader == nil {
		return &bytes.Buffer{}, nil
	}

	// the reader might already be a Buffer, in which case return a new one
	// from its underlying byte array
	if buf, ok := reader.(*bytes.Buffer); ok {
		return bytes.NewBuffer(buf.Bytes()), nil
	}

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return nil, err
	}

	TryClose(reader)

	return &buf, nil
}

// TryClose tries to cast all provided params to io.Closer, and, if
// successful, calls Close on each.
func TryClose(maybeClosers ...interface{}) {
	for _, maybeCloser := range maybeClosers {
		if closer, ok := maybeCloser.(io.Closer); ok {
			_ = closer.Close()
		}
	}
}
