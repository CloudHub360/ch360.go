package tests

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/ch360"
	mocks2 "github.com/waives/surf/ch360/mocks"
	"github.com/waives/surf/cmd/surf/commands"
	"github.com/waives/surf/test/generators"
	"testing"
)

type ListDocumentSuite struct {
	suite.Suite
	sut    *commands.ListDocumentsCmd
	client *mocks2.DocumentGetter
	ctx    context.Context
}

func (suite *ListDocumentSuite) SetupTest() {
	suite.client = new(mocks2.DocumentGetter)

	suite.sut = &commands.ListDocumentsCmd{
		Client: suite.client,
	}
	suite.ctx = context.Background()
}

func TestListDocumentSuiteRunner(t *testing.T) {
	suite.Run(t, new(ListDocumentSuite))
}

func (suite *ListDocumentSuite) TestGetAllDocuments_Execute_Calls_The_Client() {
	expectedDocuments := aListOfDocuments("charlie", "jo", "chris")
	suite.client.On("GetAll", mock.Anything).Return(expectedDocuments, nil)

	suite.sut.Execute(suite.ctx)

	suite.client.AssertCalled(suite.T(), "GetAll", suite.ctx)
}

func (suite *ListDocumentSuite) TestGetAllDocuments_Execute_Returns_An_Error_If_The_Documents_Cannot_Be_Retrieved() {
	expectedErr := errors.New("Failed")
	suite.client.On("GetAll", mock.Anything).Return(nil, expectedErr)

	actualErr := suite.sut.Execute(context.Background())

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func aListOfDocuments(ids ...string) ch360.DocumentList {
	expected := make(ch360.DocumentList, len(ids))

	for index, id := range ids {
		expected[index] = ch360.Document{
			Id:       id,
			Size:     generators.Int(),
			Sha256:   generators.String("sha"),
			FileType: generators.String("fileType"),
		}
	}

	return expected
}
