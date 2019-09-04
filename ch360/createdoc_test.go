package ch360_test

import (
	"bytes"
	"context"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/ch360/mocks"
	"github.com/CloudHub360/ch360.go/ch360/mocks/matchers"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/petergtz/pegomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreateDocumentForSuite struct {
	suite.Suite
	documentCreator *mocks.MockDocumentCreator
	documentDeleter *mocks.MockDocumentDeleter
	fileContents    *bytes.Buffer
	document        ch360.Document
	ctx             context.Context
}

func (suite *CreateDocumentForSuite) SetupTest() {
	documentId := generators.String("documentId")
	suite.document = ch360.Document{
		Id: documentId,
	}
	suite.fileContents = bytes.NewBuffer(generators.Bytes())

	suite.documentCreator = mocks.NewMockDocumentCreator(pegomock.WithT(suite.T()))
	suite.documentDeleter = mocks.NewMockDocumentDeleter(pegomock.WithT(suite.T()))
	suite.ctx = context.Background()

	pegomock.
		When(suite.documentCreator.Create(matchers.AnyContextContext(), matchers.AnyIoReader())).
		ThenReturn(suite.document, nil)
}

func TestCreateDocumentForSuiteRunner(t *testing.T) {
	suite.Run(t, new(CreateDocumentForSuite))
}

func (suite *CreateDocumentForSuite) Test_CreateDocumentFor_Creates_A_Doc_And_Then_Deletes_It() {

	ch360.CreateDocumentFor(suite.fileContents, suite.documentCreator, suite.documentDeleter,
		func(document ch360.Document) error {
			return nil
		})

	suite.documentCreator.
		VerifyWasCalledOnce().
		Create(suite.ctx, suite.fileContents)
	suite.documentDeleter.
		VerifyWasCalledOnce().
		Delete(suite.ctx, suite.document.Id)
}

func (suite *CreateDocumentForSuite) Test_CreateDocumentFor_Returns_Error_From_Function_Param() {
	expectedErr := errors.New("simulated error")

	actualErr := ch360.CreateDocumentFor(suite.fileContents, suite.documentCreator,
		suite.documentDeleter,
		func(document ch360.Document) error {
			return expectedErr
		})

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func (suite *CreateDocumentForSuite) Test_CreateDocumentFor_Deletes_The_Document_If_The_Fn_Returns_Err() {
	expectedErr := errors.New("simulated error")

	_ = ch360.CreateDocumentFor(suite.fileContents, suite.documentCreator,
		suite.documentDeleter,
		func(document ch360.Document) error {
			return expectedErr
		})

	suite.documentDeleter.
		VerifyWasCalledOnce().
		Delete(suite.ctx, suite.document.Id)
}
