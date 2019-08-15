package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/cmd/surf/services"
	mocks2 "github.com/CloudHub360/ch360.go/cmd/surf/services/mocks"
	"github.com/CloudHub360/ch360.go/pool"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

type ParallelFilesProcessorSuite struct {
	suite.Suite
	sut              *services.ParallelFilesProcessor
	progressHandler  *mocks2.ProgressHandler
	processorFactory services.ProcessorFuncFactory

	processorFunc pool.ProcessorFunc
	handlerFunc   pool.HandlerFunc

	ctx context.Context
}

func (suite *ParallelFilesProcessorSuite) SetupTest() {
	suite.processorFactory = func(ctx context.Context, filename string) pool.ProcessorFunc {
		return suite.processorFunc
	}

	suite.progressHandler = new(mocks2.ProgressHandler)

	suite.sut = &services.ParallelFilesProcessor{
		ProgressHandler: suite.progressHandler,
	}
	suite.ctx = context.Background()

	suite.processorFunc = func() (interface{}, error) {
		return nil, nil
	}
	suite.handlerFunc = func(value interface{}, e error) {
	}

	suite.progressHandler.On("NotifyStart", mock.Anything).Return(nil)
	suite.progressHandler.On("NotifyFinish").Return(nil)
	suite.progressHandler.On("Notify", mock.Anything, mock.Anything).Return(nil)
	suite.progressHandler.On("NotifyErr", mock.Anything, mock.Anything).Return(nil)

	rand.Seed(time.Now().Unix())
}

func TestParallelFilesProcessorSuiteRunner(t *testing.T) {
	suite.Run(t, new(ParallelFilesProcessorSuite))
}

func (suite *ParallelFilesProcessorSuite) Test_Err_Returned_When_Files_Glob_Matches_No_Files() {
	err := suite.sut.RunWithGlob(suite.ctx, []string{"not-present/*.pdf"}, rand.Int(),
		suite.processorFactory)

	suite.Assert().Equal(services.ErrGlobMatchesNoFiles, err)
}

func (suite *ParallelFilesProcessorSuite) Test_Err_Returned_When_File_Is_Not_Present() {

	err := suite.sut.RunWithGlob(suite.ctx, []string{"not-present.pdf"}, rand.Int(),
		suite.processorFactory)

	suite.Assert().Equal(services.ErrGlobMatchesNoFiles, err)
}

func (suite *ParallelFilesProcessorSuite) Test_ProcessorFunc_Called_Once_Per_File() {
	// Arrange
	var (
		files                = someTempFiles(2)
		expectedProcessCalls = 2
		processorFactory     = &countingProcessorFactory{}
	)
	defer deleteFiles(files)

	// Act
	suite.sut.Run(suite.ctx, files, 1, processorFactory.ProcessorFor)

	// Assert
	suite.Assert().Equal(expectedProcessCalls, processorFactory.processorCalls)
	suite.Assert().Equal(expectedProcessCalls, processorFactory.processorFactoryCalls)
}

func (suite *ParallelFilesProcessorSuite) Test_ProcessorFunc_Called_In_Parallel() {
	// Arrange
	var (
		parallelism      = 5
		delayMs          = 10 * time.Millisecond
		files            = someTempFiles(parallelism)
		processorFactory = &sleepingProcessorFactory{
			delay: delayMs,
		}

		seriesDuration = time.Duration(parallelism) * delayMs
	)
	defer deleteFiles(files)

	// Act
	start := time.Now()
	suite.sut.Run(suite.ctx, files, parallelism, processorFactory.ProcessorFor)
	parallelDuration := time.Since(start)

	// Assert
	suite.Assert().True(seriesDuration > parallelDuration)
}

func (suite *ParallelFilesProcessorSuite) Test_ProgressHandler_NotifyStart_And_NotifyFinish_Called() {
	var (
		filesCount = 5
		files      = someTempFiles(filesCount)
	)
	defer deleteFiles(files)
	suite.sut.Run(suite.ctx, files, rand.Int(), suite.processorFactory)

	suite.progressHandler.AssertCalled(suite.T(), "NotifyStart", filesCount)
	suite.progressHandler.AssertCalled(suite.T(), "NotifyFinish")
}

func (suite *ParallelFilesProcessorSuite) Test_First_Error_From_Processor_Func_Returned() {
	// Arrange
	var (
		filesCount       = 5
		files            = someTempFiles(filesCount)
		processorFactory = erroringProcessorFactory{}
		expectedErr      = processorFactory.Err(1)
	)
	defer deleteFiles(files)

	// Act
	receivedErr := suite.sut.Run(suite.ctx, files, 1, processorFactory.ProcessorFor)

	// Assert
	suite.Assert().Equal(expectedErr, receivedErr)
}

var _ services.ProcessorFuncFactory = (*countingProcessorFactory)(nil).ProcessorFor

type countingProcessorFactory struct {
	processorFactoryCalls int
	processorCalls        int
}

func (f *countingProcessorFactory) ProcessorFor(ctx context.Context, filename string) pool.ProcessorFunc {
	f.processorFactoryCalls++
	return func() (interface{}, error) {
		f.processorCalls++
		return nil, nil
	}
}

var _ services.ProcessorFuncFactory = (*erroringProcessorFactory)(nil).ProcessorFor

type erroringProcessorFactory struct {
	processorFactoryCalls int
	processorCalls        int
}

func (f *erroringProcessorFactory) ProcessorFor(ctx context.Context, filename string) pool.ProcessorFunc {
	f.processorFactoryCalls++
	return func() (interface{}, error) {
		f.processorCalls++
		return nil, f.Err(f.processorCalls)
	}
}

func (f *erroringProcessorFactory) Err(i int) error {
	return errors.New(fmt.Sprintf("Error %d", i))
}

var _ services.ProcessorFuncFactory = (*sleepingProcessorFactory)(nil).ProcessorFor

type sleepingProcessorFactory struct {
	delay time.Duration
}

func (f *sleepingProcessorFactory) ProcessorFor(ctx context.Context, filename string) pool.ProcessorFunc {
	return func() (interface{}, error) {
		time.Sleep(f.delay)
		return nil, nil
	}
}

func someTempFiles(count int) []string {
	var files []string
	for i := 0; i < count; i++ {
		f, _ := ioutil.TempFile("", "ParallelFilesProcessor_test")
		files = append(files, f.Name())
	}
	return files
}

func deleteFiles(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}
