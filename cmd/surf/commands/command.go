package commands

import (
	"context"
	"fmt"
	"os"
)

type Command interface {
	Execute(ctx context.Context) error
}

func exitOnErr(err error) {
	if err != nil && err != context.Canceled {

		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
