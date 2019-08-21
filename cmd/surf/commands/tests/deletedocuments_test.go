package tests

import (
	"context"
	"errors"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/cmd/surf/commands/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DeleteDocumentSuite struct {
	suite.Suite
	sut         *commands.DeleteDocumentCmd
	client      *mocks.DocumentDeleterGetter
	documentIds []string
	ctx         context.Context
	deleteAll   bool
}

func (suite *DeleteDocumentSuite) SetupTest() {
	suite.client = new(mocks.DocumentDeleterGetter)
	suite.documentIds = []string{"charlie", "jo", "chris"}
	suite.client.
		On("GetAll", mock.Anything).
		Return(aListOfDocuments(suite.documentIds...), nil)

	suite.client.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	suite.deleteAll = false

	suite.sut = &commands.DeleteDocumentCmd{
		Client:      suite.client,
		DocumentIDs: suite.documentIds,
		DeleteAll:   suite.deleteAll,
	}
	suite.ctx = context.Background()
}

func TestDeleteDocumentSuiteRunner(t *testing.T) {
	suite.Run(t, new(DeleteDocumentSuite))
}

func (suite *DeleteDocumentSuite) TestDeleteDocument_Execute_Deletes_The_Named_Document() {
	_ = suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
	suite.client.AssertCalled(suite.T(), "Delete", suite.ctx, "charlie")
}

func (suite *DeleteDocumentSuite) TestDeleteDocument_Execute_Returns_An_Error_If_The_Documents_Cannot_Be_Retrieved_When_Deleting_All_Documents() {
	suite.sut.DeleteAll = true
	suite.client.ExpectedCalls = nil
	expectedErr := errors.New("Failed")
	suite.client.
		On("GetAll", mock.Anything).
		Return(nil, expectedErr)

	actualErr := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expectedErr, actualErr)
	suite.client.AssertNotCalled(suite.T(), "Delete")
}

func (suite *DeleteDocumentSuite) TestDeleteDocument_Execute_Returns_An_Error_If_The_Document_Cannot_Be_Deleted() {
	suite.client.ExpectedCalls = nil
	suite.client.
		On("GetAll", mock.Anything).
		Return(aListOfDocuments(suite.documentIds...), nil)
	expectedErr := errors.New("Failed")
	suite.client.
		On("Delete", mock.Anything, mock.Anything).
		Return(expectedErr)

	actualErr := suite.sut.Execute(suite.ctx)

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func (suite *DeleteDocumentSuite) TestDeleteDocument_Retrieves_All_Documents_When_DeleteAll_Is_Specified() {
	suite.sut.DeleteAll = true

	_ = suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
}
