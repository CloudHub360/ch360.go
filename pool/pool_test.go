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
	sleepMs := 5 * time.Millisecond
	jobs := pool.MakeJobs(workerCount,
		func() (interface{}, error) {
			// Each worker will sleep for 5ms, but in parallel
			time.Sleep(sleepMs)
			return nil, nil
		},
		func(interface{}, error) {})

	// Act
	p := pool.NewPool(jobs, workerCount)
	start := time.Now()
	p.Run(context.Background())
	timeTaken := time.Since(start)

	// Assert
	assert.True(suite.T(), timeTaken < time.Duration(workerCount)*sleepMs)
}

func (suite *PoolSuite) Test_Pool_Performs_All_Jobs() {
	// Arrange
	workerCount := 10
	var jobsCompletedFlag int32 = 0
	jobs := pool.MakeJobs(workerCount,
		func() (interface{}, error) {
			// Executed in parallel
			atomic.AddInt32(&jobsCompletedFlag, 1)
			return nil, nil
		},
		func(interface{}, error) {})

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
		func() (interface{}, error) {
			return expectedResult, expectedError
		},
		func(result interface{}, err error) {
			receivedResult = result.(string)
			receivedErr = err
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
		jobsRun          int
		jobsCount        = 10
		allowedJobsCount = 2
		ctx, cancel      = context.WithCancel(context.Background())
	)
	jobs := pool.MakeJobs(jobsCount,
		func() (interface{}, error) {
			jobsRun++
			if jobsRun == allowedJobsCount {
				cancel()
			}

			return nil, nil
		},
		func(interface{}, error) {})

	// Act
	p := pool.NewPool(jobs, 1)
	p.Run(ctx)

	// Assert
	// may be another job queued, so we allow an extra to be processed
	// after cancel
	assert.True(suite.T(), jobsRun <= allowedJobsCount+1)
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
		func() (interface{}, error) {
			cancel()

			return nil, nil
		},
		func(interface{}, error) {
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
		func() (interface{}, error) {
			return nil, context.Canceled
		},
		func(interface{}, error) {
			handlersRun++
		})

	// Act
	p := pool.NewPool(jobs, 1)
	p.Run(context.Background())

	// Assert
	assert.Equal(suite.T(), 0, handlersRun)
}
