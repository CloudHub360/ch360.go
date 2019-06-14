package ch360

import (
	"bytes"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/net"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type LoggingDoer struct {
	wrappedSender net.HttpDoer
	out           io.Writer
}

func NewLoggingDoer(httpDoer net.HttpDoer, out io.Writer) *LoggingDoer {
	return &LoggingDoer{
		wrappedSender: httpDoer,
		out:           out,
	}
}

func (d *LoggingDoer) Do(request *http.Request) (*http.Response, error) {
	requestBytes, err := httputil.DumpRequestOut(request, false)

	if err != nil {
		return nil, err
	}
	_, err = d.out.Write(requestBytes)

	if err != nil {
		return nil, err
	}

	if request.Body != nil {
		body, err := request.GetBody()

		if err != nil {
			return nil, err
		}

		bodyBuffer := bytes.Buffer{}
		_, _ = bodyBuffer.ReadFrom(body)

		if json.Valid(bodyBuffer.Bytes()) {
			formattedJson := bytes.Buffer{}
			_ = json.Indent(&formattedJson, bodyBuffer.Bytes(), "", "  ")
			formattedJson.WriteTo(d.out)
			d.out.Write([]byte("\n"))
		} else {
			d.out.Write([]byte("<binary request body>\n\n"))
		}
	}

	response, capturedErr := d.wrappedSender.Do(request)

	responseBodyBuffer := bytes.Buffer{}
	responseBodyBuffer.ReadFrom(response.Body)
	bodyReader := bytes.NewReader(responseBodyBuffer.Bytes())
	response.Body = ioutil.NopCloser(bodyReader)

	responseBytes, err := httputil.DumpResponse(response, false)

	d.out.Write(responseBytes)

	if json.Valid(responseBodyBuffer.Bytes()) {
		formattedJson := bytes.Buffer{}
		json.Indent(&formattedJson, responseBodyBuffer.Bytes(), "", "  ")
		formattedJson.WriteTo(d.out)
	} else {
		d.out.Write([]byte("<binary response body>\n\n"))
	}

	d.out.Write([]byte("\n"))
	return response, capturedErr
}
