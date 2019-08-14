package commands

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/mattn/go-zglob"
)

var ErrGlobMatchesNoFiles = errors.New("file pattern does not match any files")

type ProgressHandler interface {
	Notify(filename string, result interface{}) error
	NotifyErr(filename string, err error) error
	NotifyStart(totalJobs int) error
	NotifyFinish() error
}

type ParallelFilesProcessor struct {
	ProgressHandler ProgressHandler
}

//go:generate mockery -name "ProcessorFuncFactory"
type ProcessorFuncFactory func(ctx context.Context, filename string) pool.ProcessorFunc

func (p *ParallelFilesProcessor) RunWithGlob(ctx context.Context,
	filesPatterns []string,
	parallelism int,
	processorFuncFactory ProcessorFuncFactory) error {

	var files []string
	for _, filePattern := range filesPatterns {
		globResult, _ := zglob.Glob(filePattern)
		files = append(files, globResult...)
	}
	fileCount := len(files)
	if fileCount == 0 {
		return ErrGlobMatchesNoFiles
	}

	return p.Run(ctx, files, parallelism, processorFuncFactory)
}

func (p *ParallelFilesProcessor) Run(ctx context.Context,
	files []string,
	parallelism int,
	processorFuncFactory ProcessorFuncFactory) error {

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
			processorFuncFactory(ctx, filename),
			func(result interface{}, e error) {
				if e != nil {
					errs = append(errs, e)
					p.ProgressHandler.NotifyErr(filename, e)
					cancel()
				} else {
					if e = p.ProgressHandler.Notify(filename, result); e != nil {
						// An error occurred while writing output
						errs = append(errs, e)
						cancel()
					}
				}
			})

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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
