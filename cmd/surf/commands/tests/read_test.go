package tests

import (
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type readCommandSuite struct {
	suite.Suite
	readerService *mocks.ReaderService
	sut           *commands.ReadCmd
	filePatterns  []string
	ctx           context.Context
	readMode      ch360.ReadMode
	expectedErr   error
}

func (suite *readCommandSuite) SetupTest() {
	suite.readerService = new(mocks.ReaderService)
	suite.ctx, _ = context.WithCancel(context.Background())
	suite.filePatterns = []string{generators.String("file1"), generators.String("file2")}

	suite.expectedErr = errors.New("simulated error")

	suite.readerService.
		On("ReadAll", mock.Anything, mock.Anything, mock.Anything).
		Return(suite.expectedErr)

	suite.readMode = ch360.ReadPDF

	suite.sut = &commands.ReadCmd{
		FilePaths:     suite.filePatterns,
		ReaderService: suite.readerService,
		ReadMode:      suite.readMode,
	}
}

func TestReadCommandSuiteRunner(t *testing.T) {
	suite.Run(t, new(readCommandSuite))
}

func (suite *readCommandSuite) Test_ReaderService_ReadAll_Called_With_Correct_Params() {
	_ = suite.sut.Execute(suite.ctx)

	suite.readerService.AssertCalled(suite.T(),
		"ReadAll",
		suite.ctx,
		suite.filePatterns,
		suite.readMode)
}

func (suite *readCommandSuite) Test_Error_Returned_From_ReaderService() {
	actualErr := suite.sut.Execute(suite.ctx)

	assert.EqualError(suite.T(), errors.Cause(actualErr), suite.expectedErr.Error())
}
