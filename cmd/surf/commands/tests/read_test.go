package tests

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	ch360mocks "github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type readCommandSuite struct {
	suite.Suite
	filesProcessor  *mocks.FilesProcessor
	fileReader      *mocks.FileReader
	documentGetter  *ch360mocks.DocumentGetter
	sut             *commands.Read
	filesPattern    string
	parallelWorkers int
	ctx             context.Context
	readMode        ch360.ReadMode
}

func (suite *readCommandSuite) SetupTest() {
	suite.fileReader = new(mocks.FileReader)
	suite.filesProcessor = new(mocks.FilesProcessor)
	suite.documentGetter = new(ch360mocks.DocumentGetter)
	suite.ctx, _ = context.WithCancel(context.Background())
	suite.filesPattern = generators.String("filesPattern")
	suite.parallelWorkers = 10

	suite.readMode = ch360.ReadPDF

	suite.sut = commands.NewReadFilesCommand(suite.fileReader,
		suite.filesProcessor,
		suite.readMode,
		suite.filesPattern, suite.parallelWorkers, suite.documentGetter)

	// set up the happy path
	suite.documentGetter.
		On("GetAll", mock.Anything).
		Return(nil, nil)

	suite.filesProcessor.
		On("RunWithGlob", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	suite.fileReader.
		On("Read", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
}

func TestReadCommandSuiteRunner(t *testing.T) {
	suite.Run(t, new(readCommandSuite))
}

func (suite *readCommandSuite) Test_FileProcessor_RunWithGlob_Called_With_Correct_Params() {
	suite.sut.Execute(suite.ctx)

	suite.filesProcessor.AssertCalled(suite.T(),
		"RunWithGlob",
		suite.ctx,
		suite.filesPattern,
		suite.parallelWorkers,
		suite.sut)
}

func (suite *readCommandSuite) Test_Process_Factory_Returns_Func_Which_Calls_FileReader() {
	// pass the name of a file that exists (this one)
	processorFn := suite.sut.ProcessorFor(suite.ctx, "read_test.go")

	processorFn()

	suite.fileReader.AssertCalled(suite.T(),
		"Read",
		suite.ctx,
		mock.Anything,
		suite.readMode)
}

func (suite *readCommandSuite) Test_Process_Factory_Returns_Err_If_File_Does_Not_Exist() {
	// pass the name of a file that doesn't exist
	processorFn := suite.sut.ProcessorFor(suite.ctx, "does-not-exist")

	_, err := processorFn()

	suite.Assert().True(strings.HasPrefix(err.Error(), "Error reading file does-not-exist"))
}

func (suite *readCommandSuite) Test_DocumentGetter_Called_To_Calculate_Parallelism() {
	var testData = []struct {
		docsInAccount       []ch360.Document
		totalDocSlots       int
		expectedParalellism int
	}{
		{
			docsInAccount:       make([]ch360.Document, 10),
			totalDocSlots:       20,
			expectedParalellism: 10,
		},
		{
			docsInAccount:       nil,
			totalDocSlots:       10,
			expectedParalellism: 10,
		},
		{
			docsInAccount:       nil,
			totalDocSlots:       5,
			expectedParalellism: 5,
		},
	}

	origTotalDocSlots := ch360.TotalDocumentSlots
	defer func() {
		ch360.TotalDocumentSlots = origTotalDocSlots
	}()

	for _, td := range testData {
		suite.documentGetter.ExpectedCalls = nil
		suite.documentGetter.
			On("GetAll", mock.Anything).
			Return(td.docsInAccount, nil)
		ch360.TotalDocumentSlots = td.totalDocSlots

		suite.sut.Execute(suite.ctx)

		suite.filesProcessor.AssertCalled(suite.T(), "RunWithGlob",
			suite.ctx,
			suite.filesPattern,
			td.expectedParalellism,
			suite.sut)
	}
}
