package commands

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/mattn/go-zglob"
)

var ErrGlobMatchesNoFiles = errors.New("file pattern does not match any files")

type ParallelFilesProcessor struct {
	ProgressHandler ProgressHandler
}

//go:generate mockery -name "ProcessorFuncFactory"
type ProcessorFuncFactory interface {
	ProcessorFor(ctx context.Context, filename string) pool.ProcessorFunc
}

//go:generate mockery -name "HandlerFuncFactory"
type HandlerFuncFactory interface {
	HandlerFor(cancel context.CancelFunc, filename string, progressHandler ProgressHandler, errs *[]error) pool.HandlerFunc
}

func (p *ParallelFilesProcessor) RunWithGlob(ctx context.Context,
	filesPattern string,
	parallelism int,
	processorFuncFactory ProcessorFuncFactory,
	handlerFuncFactory HandlerFuncFactory) error {

	files, _ := zglob.Glob(filesPattern)
	fileCount := len(files)
	if fileCount == 0 {
		return ErrGlobMatchesNoFiles
	}

	return p.Run(ctx, files, parallelism, processorFuncFactory, handlerFuncFactory)
}

func (p *ParallelFilesProcessor) Run(ctx context.Context,
	files []string,
	parallelism int,
	processorFuncFactory ProcessorFuncFactory,
	handlerFuncFactory HandlerFuncFactory) error {

	ctx, cancel := context.WithCancel(ctx)

	var (
		processFileJobs []pool.Job
		errs            []error
	)

	for _, filename := range files {
		// The memory of the 'filename' var is reused here, see:
		// https://golang.org/doc/faq#closures_and_goroutines
		// The workaround is to copy it:
		filename := filename // <- copy

		processFileJob := pool.NewJob(
			processorFuncFactory.ProcessorFor(ctx, filename),
			handlerFuncFactory.HandlerFor(cancel, filename, p.ProgressHandler, &errs))

		processFileJobs = append(processFileJobs, processFileJob)
	}

	workPool := pool.NewPool(processFileJobs, min(parallelism, len(files)))

	p.ProgressHandler.NotifyStart(len(processFileJobs))
	defer p.ProgressHandler.NotifyFinish()
	workPool.Run(ctx)

	// Just return the first error.
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
