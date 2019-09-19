package net

import (
	"bytes"
	"github.com/waives/surf/ioutils"
	"io/ioutil"
	"net/http"
)

// RequestBodyBytes reads the body from a request into a bytes.Buffer,
// resets the request Body and returns the buffer.
func RequestBodyBytes(request *http.Request) *bytes.Buffer {
	bodyBuffer, err := ioutils.DrainClose(request.Body)

	if err != nil {
		return nil
	}

	request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer.Bytes()))
	return bodyBuffer
}
