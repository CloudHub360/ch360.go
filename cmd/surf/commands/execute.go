package commands

import (
	"fmt"
	"io"
	"strings"
)

// Execute prints the provided message to out, runs fn(), then
// prints [OK] or [FAILED] depending on fn's success.
func ExecuteWithMessage(message string, out io.Writer, fn func() error) error {
	if !strings.HasSuffix(message, " ") {
		message += " "
	}

	_, err := fmt.Fprint(out, message)

	if err = fn(); err != nil {
		_, _ = fmt.Fprintln(out, "[FAILED]")
		return err
	}

	_, _ = fmt.Fprintln(out, "[OK]")
	return nil
}
