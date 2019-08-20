package commands

import (
	"fmt"
	"os"
	"strings"
)

// Execute prints the provided message to stderr, runs fn(), then
// prints [OK] or [FAILED] depending on fn's success.
func ExecuteWithMessage(message string, fn func() error) error {
	if !strings.HasSuffix(message, " ") {
		message += " "
	}

	out := os.Stderr

	_, err := fmt.Fprint(out, message)

	if err = fn(); err != nil {
		_, _ = fmt.Fprintln(out, "[FAILED]")
		return err
	}

	_, _ = fmt.Fprintln(out, "[OK]")
	return nil
}
