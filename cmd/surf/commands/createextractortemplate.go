package commands

import (
	"context"
	"io"
)

const CreateExtractorTemplateCommand = "create extractor-template"

var _ Command = (*CreateExtractorTemplate)(nil)

type CreateExtractorTemplate struct {
	out io.Writer
}

func (cmd CreateExtractorTemplate) Execute(ctx context.Context) error {
	return Execute("Creating extractor template...", cmd.out, func() error {
		return nil
	})
}

func (cmd CreateExtractorTemplate) Usage() string {
	return CreateExtractorTemplateCommand
}
