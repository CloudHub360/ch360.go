package io_util

import "io"

// AutoCloser wraps an io.WriteCloser and calls Close() after any call to Write().
type AutoCloser struct {
	Underlying io.WriteCloser
}

func (ac *AutoCloser) Close() error {
	return ac.Underlying.Close()
}

func (ac *AutoCloser) Write(p []byte) (n int, err error) {
	n, err = ac.Underlying.Write(p)

	closeErr := ac.Underlying.Close()

	if err == nil {
		err = closeErr
	}

	return
}
