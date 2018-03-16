package formatters

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

var _ ResultsFormatter = (*NoopResultsFormatter)(nil)

type NoopResultsFormatter struct {
}

func NewReadResultsFormatter() *NoopResultsFormatter {
	return &NoopResultsFormatter{}
}

func ErrUnexpectedType(val interface{}) error {
	return errors.New(fmt.Sprintf("Unexpected type: %T", val))
}

func (f *NoopResultsFormatter) WriteResult(writer io.Writer, fullPath string, result interface{}, options FormatOption) error {

	readCloser, ok := result.(io.ReadCloser)

	if !ok {
		return ErrUnexpectedType(result)
	}
	defer readCloser.Close()

	_, err := io.Copy(writer, readCloser)
	return err
}

func (f *NoopResultsFormatter) Flush(writer io.Writer) error {
	return nil
}
