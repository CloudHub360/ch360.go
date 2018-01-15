package pool_test

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync/atomic"
	"testing"
	"time"
)

type PoolSuite struct {
	suite.Suite
	sut pool.Pool
}

func (suite *PoolSuite) SetupTest() {
	suite.sut = pool.NewPool(nil, 1)
}

func TestPoolSuiteRunner(t *testing.T) {
	suite.Run(t, new(PoolSuite))
}

func (suite *PoolSuite) Test_Pool_Performs_Work_In_Parallel() {
	// Arrange
	workerCount := 10
	sleepMs := 5
	jobs := pool.MakeJobs(workerCount,
		func() pool.JobResult {
			// Each worker will sleep for 5ms, but in parallel
			time.Sleep(time.Millisecond * time.Duration(sleepMs))
			return pool.JobResult{}
		},
		func(result pool.JobResult) {})

	// Act
	p := pool.NewPool(jobs, workerCount)
	start := time.Now()
	p.Run(context.Background())
	timeTaken := time.Since(start)

	// Assert
	assert.True(suite.T(), timeTaken < time.Duration(workerCount*sleepMs)*time.Millisecond)
}

func (suite *PoolSuite) Test_Pool_Performs_All_Jobs() {
	// Arrange
	workerCount := 10
	var jobsCompletedFlag int32 = 0
	jobs := pool.MakeJobs(workerCount,
		func() pool.JobResult {
			// Executed in parallel
			atomic.AddInt32(&jobsCompletedFlag, 1)
			return pool.JobResult{}
		},
		func(result pool.JobResult) {})

	// Act
	p := pool.NewPool(jobs, workerCount)
	p.Run(context.Background())

	// Assert
	assert.Equal(suite.T(), int32(workerCount), jobsCompletedFlag)
}

func (suite *PoolSuite) Test_Pool_Calls_Handler_With_JobResults() {
	// Arrange
	var (
		expectedResult = generators.String("pool")
		expectedError  = errors.New("err")
		receivedResult string
		receivedErr    error
	)

	jobs := pool.MakeJobs(1,
		func() pool.JobResult {
			return pool.JobResult{
				Value: expectedResult,
				Err:   expectedError,
			}
		},
		func(result pool.JobResult) {
			receivedResult = result.Value.(string)
			receivedErr = result.Err
		})

	// Act
	p := pool.NewPool(jobs, 1)
	p.Run(context.Background())

	// Assert
	assert.Equal(suite.T(), expectedResult, receivedResult)
	assert.Equal(suite.T(), expectedError, receivedErr)
}

// test context cancellation prevents more jobs from being run
func (suite *PoolSuite) Test_Pool_Does_Not_Process_Jobs_After_Context_Cancel() {
	// Arrange
	var (
		jobsRun          int32
		jobsCount        = int(10)
		allowedJobsCount = 2
		ctx, cancel      = context.WithCancel(context.Background())
		jobCompleteSig   = make(chan bool, 0)
	)
	jobs := pool.MakeJobs(jobsCount,
		func() pool.JobResult {
			jobCompleteSig <- true

			time.Sleep(time.Millisecond)
			atomic.AddInt32(&jobsRun, 1)

			return pool.JobResult{}
		},
		func(result pool.JobResult) {})

	// Act
	go func() {
		// Cancel the context after 2 jobs
		for _ = range jobCompleteSig {
			if jobsCount >= allowedJobsCount {
				cancel()
			}
		}
	}()

	p := pool.NewPool(jobs, 1)
	p.Run(ctx)
	close(jobCompleteSig)

	// Assert
	assert.Equal(suite.T(), int32(allowedJobsCount), jobsRun)
}

// test context cancellation prevents results handlers from being called
func (suite *PoolSuite) Test_Pool_Does_Not_Process_Handler_After_Context_Cancel() {
	// Arrange
	var (
		jobsCount   = 10
		handlersRun int
		ctx, cancel = context.WithCancel(context.Background())
	)
	jobs := pool.MakeJobs(jobsCount,
		func() pool.JobResult {
			cancel()

			return pool.JobResult{}
		},
		func(result pool.JobResult) {
			handlersRun++
		})

	// Act
	p := pool.NewPool(jobs, 1)
	p.Run(ctx)

	// Assert
	assert.True(suite.T(), handlersRun < jobsCount)
}

func (suite *PoolSuite) Test_Pool_Does_Not_Process_Handler_After_Job_Returns_Context_Cancelled_Error() {
	// Arrange
	var (
		jobsCount   = 10
		handlersRun int
	)
	jobs := pool.MakeJobs(jobsCount,
		func() pool.JobResult {
			return pool.JobResult{
				Err: context.Canceled,
			}
		},
		func(result pool.JobResult) {
			handlersRun++
		})

	// Act
	p := pool.NewPool(jobs, 1)
	p.Run(context.Background())

	// Assert
	assert.Equal(suite.T(), 0, handlersRun)
}
