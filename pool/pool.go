package pool

import (
	"context"
	"sync"
)

type Pool struct {
	jobs    []Job
	workers int
}

type Job struct {
	process ProcessorFunc
	handle  HandlerFunc
}

type ProcessorFunc func() (interface{}, error)
type HandlerFunc func(interface{}, error)

type jobResult struct {
	Err   error
	Value interface{}
}

type jobAndResult struct {
	job    Job
	result jobResult
}

func NewPool(jobs []Job, workers int) Pool {
	return Pool{
		jobs:    jobs,
		workers: workers,
	}
}

func NewJob(process ProcessorFunc, handler HandlerFunc) Job {
	return Job{
		process: process,
		handle:  handler,
	}
}

func MakeJobs(n int, process ProcessorFunc, handle HandlerFunc) []Job {
	var jobs []Job

	for i := 0; i < n; i++ {
		job := NewJob(process, handle)
		jobs = append(jobs, job)
	}

	return jobs
}

func (p *Pool) sourceJobs(ctx context.Context, jobsChan chan Job) {
	defer close(jobsChan)

	for _, job := range p.jobs {
		select {
		case <-ctx.Done():
			return
		default:
			jobsChan <- job
		}
	}
}

func (p *Pool) processJob(ctx context.Context, jobsChan chan Job, resultsChan chan jobAndResult, wg *sync.WaitGroup) {
	for job := range jobsChan {
		result, err := job.process()

		resultsChan <- jobAndResult{
			result: jobResult{
				Err:   err,
				Value: result,
			},
			job: job,
		}
	}
	wg.Done()
}

func (p *Pool) handleResult(jobRes *jobAndResult) {
	job, result := jobRes.job, jobRes.result

	if result.Err == context.Canceled {
		return
	}

	job.handle(result.Value, result.Err)
}

func (p *Pool) Run(ctx context.Context) {
	// Set up jobs channel
	jobsChan := make(chan Job)

	// Begin adding jobs
	go p.sourceJobs(ctx, jobsChan)

	resultsChan := make(chan jobAndResult)

	// Start processing in parallel
	wg := sync.WaitGroup{}
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go p.processJob(ctx, jobsChan, resultsChan, &wg)
	}

	go func() {
		// Wait for all workers to complete, then close the results channel
		wg.Wait()
		close(resultsChan)
	}()

	// Handle results in calling goroutine
	for jobRes := range resultsChan {
		p.handleResult(&jobRes)
	}
}
