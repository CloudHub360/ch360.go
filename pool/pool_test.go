package pool_test

import (
	"context"
	"github.com/CloudHub360/ch360.go/pool"
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

// test all jobs are performed
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

// test results handlers are called with results

// test context cancellation prevents more jobs from being run

// test context cancellation prevents results handlers from being called
