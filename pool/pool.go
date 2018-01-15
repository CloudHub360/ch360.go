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
	do     ProcessorFunc
	handle HandlerFunc
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

func NewJob(worker ProcessorFunc, handler HandlerFunc) Job {
	return Job{
		do:     worker,
		handle: handler,
	}
}

func MakeJobs(n int, do ProcessorFunc, handle HandlerFunc) []Job {
	var jobs []Job

	for i := 0; i < n; i++ {
		job := NewJob(do, handle)
		jobs = append(jobs, job)
	}

	return jobs
}

func (p *Pool) Run(ctx context.Context) {
	// Set up jobs channel
	jobsChan := make(chan Job)
	go func() {
		defer close(jobsChan)

		for _, job := range p.jobs {
			select {
			case <-ctx.Done():
				return
			default:
				jobsChan <- job
			}
		}
	}()

	// Set up results channel
	resultsChan := make(chan jobAndResult)

	// Start processing in background
	wg := sync.WaitGroup{}
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			for job := range jobsChan {
				result, err := job.do()

				resultsChan <- jobAndResult{
					result: jobResult{
						Err:   err,
						Value: result,
					},
					job: job,
				}
			}
			wg.Done()
		}()
	}

	go func() {
		// Wait for all workers to complete, then close the results channel
		wg.Wait()
		close(resultsChan)
	}()

	// Handle results in calling goroutine
	for jobRes := range resultsChan {
		job, result := jobRes.job, jobRes.result

		if result.Err == context.Canceled {
			continue
		}

		job.handle(result.Value, result.Err)
	}
}
