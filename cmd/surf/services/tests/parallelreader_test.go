package tests

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/ch360"
	ch360mocks "github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/cmd/surf/services"
	"github.com/CloudHub360/ch360.go/cmd/surf/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type parallelReaderSuite struct {
	suite.Suite
	fileReader      *mocks.FileReader
	progressHandler *mocks.ProgressHandler
	documentGetter  *ch360mocks.DocumentGetter
	sut             *services.ParallelReaderService
	filePatterns    []string
	ctx             context.Context
	readMode        ch360.ReadMode
}

func (suite *parallelReaderSuite) SetupTest() {
	suite.fileReader = new(mocks.FileReader)
	suite.documentGetter = new(ch360mocks.DocumentGetter)
	suite.progressHandler = new(mocks.ProgressHandler)
	suite.ctx, _ = context.WithCancel(context.Background())
	suite.filePatterns = []string{"testdata/empty-file1.txt", "testdata/empty-file2.txt"}

	suite.readMode = ch360.ReadPDF

	suite.sut = services.NewParallelReaderService(suite.fileReader, suite.documentGetter, suite.progressHandler)

	suite.documentGetter.
		On("GetAll", mock.Anything).
		Return(nil, nil)

	suite.progressHandler.
		On("NotifyStart", mock.Anything).
		Return(nil)

	suite.progressHandler.
		On("Notify", mock.Anything, mock.Anything).
		Return(nil)

	suite.progressHandler.
		On("NotifyFinish").
		Return(nil)

	suite.progressHandler.
		On("NotifyErr", mock.Anything, mock.Anything).
		Return(nil)

	suite.fileReader.
		On("Read", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
}

func TestReadCommandSuiteRunner(t *testing.T) {
	suite.Run(t, new(parallelReaderSuite))
}

func (suite *parallelReaderSuite) Test_ProgressHandler_Notify_Called_Once_Per_File() {
	_ = suite.sut.ReadAll(suite.ctx, suite.filePatterns, suite.readMode)

	for _, filename := range suite.filePatterns {
		suite.progressHandler.AssertCalled(suite.T(),
			"Notify",
			filename,
			nil)
	}
}

func (suite *parallelReaderSuite) Test_ReadAll_Returns_Error_If_File_Does_Not_Exist() {
	filename := "does-not-exist.txt"

	err := suite.sut.ReadAll(suite.ctx, []string{filename}, suite.readMode)

	assert.Error(suite.T(), err)
}

func (suite *parallelReaderSuite) Test_DocumentGetter_Called_To_Calculate_Parallelism() {
	_ = suite.sut.ReadAll(suite.ctx, suite.filePatterns, suite.readMode)

	suite.documentGetter.AssertCalled(suite.T(),
		"GetAll", suite.ctx)
}

func (suite *parallelReaderSuite) Test_ReadAll_Returns_Error_If_ExtractDocument_Fails() {
	suite.fileReader.ExpectedCalls = nil
	readErr := errors.New("simulated error")

	suite.fileReader.
		On("Read", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, readErr)

	err := suite.sut.ReadAll(suite.ctx, suite.filePatterns, suite.readMode)

	assert.Error(suite.T(), err)
}
