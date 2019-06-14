package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"io"
	"os"
)

//go:generate mockery -name "ExtractorCreator"

const UploadExtractorCommand = "upload extractor"

type ExtractorCreator interface {
	Create(ctx context.Context, name string, config io.Reader) error
	CreateFromJson(ctx context.Context, name string, jsonTemplate io.Reader) error
	CreateFromModules(ctx context.Context, name string, modules ch360.ModulesTemplate) error
}

type UploadExtractor struct {
	writer        io.Writer
	creator       ExtractorCreator
	extractorName string
	config        io.Reader
}

func NewUploadExtractor(writer io.Writer,
	creator ExtractorCreator,
	extractorName string,
	config io.Reader) *UploadExtractor {
	return &UploadExtractor{
		writer:        writer,
		creator:       creator,
		config:        config,
		extractorName: extractorName,
	}
}

func NewUploadExtractorFromArgs(params *config.RunParams,
	client ExtractorCreator, out io.Writer) (*UploadExtractor, error) {

	configFile, err := os.Open(params.ConfigPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The file '%s' could not be found.", params.ConfigPath))
	}

	return NewUploadExtractor(out,
		client,
		params.Name,
		configFile), nil
}

func (cmd *UploadExtractor) Execute(ctx context.Context) error {
	fmt.Fprintf(cmd.writer, "Uploading extractor '%s'... ", cmd.extractorName)

	err := cmd.creator.Create(ctx, cmd.extractorName, cmd.config)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	fmt.Fprintln(cmd.writer, "[OK]")

	return nil
}

func (cmd UploadExtractor) Usage() string {
	return UploadExtractorCommand
}
