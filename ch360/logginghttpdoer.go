package ch360

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CloudHub360/ch360.go/net"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"sync"
	"sync/atomic"
)

// LoggingDoer is an HttpDoer decorator that logs all HTTP requests and
// responses to the specified io.Writer. It indents any json request /
// response bodies, and redacts any non-json bodies.
type LoggingDoer struct {
	wrappedSender net.HttpDoer
	out           io.Writer
	mutex         sync.Mutex
	count         uint32
}

func NewLoggingDoer(httpDoer net.HttpDoer, out io.Writer) *LoggingDoer {
	return &LoggingDoer{
		wrappedSender: httpDoer,
		out:           out,
	}
}

func (d *LoggingDoer) Do(request *http.Request) (*http.Response, error) {
	requestId := atomic.AddUint32(&d.count, 1)
	requestBytes := d.formatRequest(request, requestId)

	d.safeWrite(requestBytes)

	response, capturedErr := d.wrappedSender.Do(request)

	if response != nil {
		responseBytes := d.formatResponse(response, requestId)

		d.safeWrite(responseBytes)
	}

	return response, capturedErr
}

func (d *LoggingDoer) safeWrite(bytes []byte) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, _ = d.out.Write(bytes)
}

func (d *LoggingDoer) formatRequest(request *http.Request, requestId uint32) []byte {
	requestBytes, err := httputil.DumpRequestOut(request, false)

	if err != nil {
		return nil
	}

	logBuffer := bytes.NewBufferString(fmt.Sprintf("[%04d -->] ", requestId))
	logBuffer.Write(requestBytes)

	if request.Body != nil {
		body, err := request.GetBody()

		if err != nil {
			return nil
		}

		bodyBuffer := bytes.Buffer{}
		_, _ = bodyBuffer.ReadFrom(body)

		if json.Valid(bodyBuffer.Bytes()) {
			_ = json.Indent(logBuffer, bodyBuffer.Bytes(), "", "  ")
			logBuffer.WriteString("\n")
		} else {
			logBuffer.WriteString("<binary request body>\n")
		}
	}

	logBuffer.WriteString("\n")

	return logBuffer.Bytes()
}

func (d *LoggingDoer) formatResponse(response *http.Response, requestId uint32) []byte {
	// get headers
	responseHeaders, err := httputil.DumpResponse(response, false)

	if err != nil {
		return nil
	}

	logBuffer := bytes.NewBufferString(fmt.Sprintf("[%04d <--] ", requestId))

	logBuffer.Write(responseHeaders)

	bodyBuffer := bytes.Buffer{}

	// reset the response body
	if response.Body != nil {
		_, err := bodyBuffer.ReadFrom(response.Body)

		if err != nil {
			return nil
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(bodyBuffer.Bytes()))
	}

	// format json
	if json.Valid(bodyBuffer.Bytes()) {
		formattedJson := bytes.Buffer{}
		err = json.Indent(&formattedJson, bodyBuffer.Bytes(), "", "  ")

		if err != nil {
			return nil
		}

		logBuffer.Write(formattedJson.Bytes())
		logBuffer.WriteString("\n\n")
	} else {
		logBuffer.WriteString("<binary response body>\n\n")
	}

	return logBuffer.Bytes()
}
